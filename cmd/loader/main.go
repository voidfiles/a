package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

func getLcFiles() []string {
	return []string{
		"http://id.loc.gov/static/data/authoritieschildrensSubjects.nt.zip",
		"http://id.loc.gov/static/data/authoritiesdemographicTerms.nt.skos.zip",
		"http://id.loc.gov/static/data/authoritiesgenreForms.nt.zip",
		"http://id.loc.gov/static/data/authoritiesnames.nt.skos.zip",
		"http://id.loc.gov/static/data/authoritiesperformanceMediums.nt.zip",
		"http://id.loc.gov/static/data/authoritiessubjects.nt.skos.zip",
		"http://id.loc.gov/static/data/vocabularycountries.nt.zip",
		"http://id.loc.gov/static/data/vocabularyethnographicTerms.nt.zip",
		"http://id.loc.gov/static/data/vocabularygeographicAreas.nt.zip",
		"http://id.loc.gov/static/data/vocabularyiso639-1.nt.zip",
		"http://id.loc.gov/static/data/vocabularyiso639-2.nt.zip",
		"http://id.loc.gov/static/data/vocabularyiso639-5.nt.zip",
		"http://id.loc.gov/static/data/vocabularylanguages.nt.zip",
		"http://id.loc.gov/static/data/vocabularyorganizations.nt.zip",
		"http://id.loc.gov/static/data/vocabularyrelators.nt.zip",
	}
}

func downloadURL(url string, ret chan string) {
	urlParts := strings.Split(url, "/")
	filename := urlParts[len(urlParts)-1]
	tmpfile, err := ioutil.TempFile("", filename)

	defer tmpfile.Close()

	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ret <- tmpfile.Name()
}

func consumeChannel(c chan string) []string {
	s := make([]string, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}

func downloadAuthorityArchives() []string {
	var wg sync.WaitGroup
	downloadedFilePaths := make(chan string)
	urlsToDownload := getLcFiles()
	wg.Add(len(urlsToDownload))
	for _, url := range urlsToDownload {
		go func(url string, ret chan string) {
			defer wg.Done()
			downloadURL(url, ret)
		}(url, downloadedFilePaths)
	}

	wg.Wait()

	return consumeChannel(downloadedFilePaths)
}

func loadAuthorityArchiveIntoCayley(archivePaths []string) error {

	for _, path := range archivePaths {
		command := "zcat" + path + " | cayley load -c conf.yml -i -"
		out, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			return err
		}
		log.Print(out)
	}

	return nil
}

func main() {
	archivePaths := downloadAuthorityArchives()
	loadAuthorityArchiveIntoCayley(archivePaths)
}
