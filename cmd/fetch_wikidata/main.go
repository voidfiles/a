package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func fetchQuery(q string) (string, error) {
	form := url.Values{}
	form.Set("query", q)
	b := form.Encode()
	log.Printf("Yoyoyoyo here q: %v", q)
	// TODO make optional GET or Post, Query() should default GET (idempotent, cacheable)
	// maybe new for updates: func (r *Repo) Update(q string) using POST?
	req, err := http.NewRequest(
		"POST",
		"https://query.wikidata.org/sparql",
		bytes.NewBufferString(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))
	req.Header.Set("Accept", "application/sparql-results+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(rawBody), nil

}

const query = `SELECT
  ?item ?itemLabel ?p ?o
WHERE {
  VALUES (?p) {
    (wdt:P244) (wdt:P214) (wdt:P4801) (wdt:P1014) (wdt:P486)
  }
  ?item ?p ?o. SERVICE wikibase:label {
    bd:serviceParam wikibase:language '[AUTO_LANGUAGE],fr,ar,be,bg,bn,ca,cs,da,de,el,en,es,et,fa,fi,he,hi,hu,hy,id,it,ja,jv,ko,nb,nl,eo,pa,pl,pt,ro,ru,sh,sk,sr,sv,sw,te,th,tr,uk,yue,vec,vi,zh'.
  }
}`

func main() {
	// repo, err := sparql.NewRepo("https://query.wikidata.org/sparql")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	res, err := fetchQuery(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(res)
}
