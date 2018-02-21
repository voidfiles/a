package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/boutros/marc"
	"github.com/coreos/bbolt"
	"github.com/voidfiles/a/cli"
	"github.com/voidfiles/a/marcdex"
	"github.com/voidfiles/a/recordstore"
)

func detectFormat(f *os.File) (marc.Format, error) {
	sniff := make([]byte, 64)
	_, err := f.Read(sniff)
	if err != nil {
		log.Fatal(err)
	}
	format := marc.DetectFormat(sniff)

	// rewind reader
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	switch format {
	case marc.MARC, marc.LineMARC, marc.MARCXML:
		return format, nil
	default:
		return format, errors.New("unknown MARC format")
	}
}

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	args := cli.GetArgs()

	db, err := bolt.Open(args.Dbpath, 600, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	if err != nil {
		panic(err)
	}
	recordStore := recordstore.MustNewRecordStore(db)
	var reader io.Reader
	var format marc.Format
	if args.InputPath != "" {
		reader, err = os.Open(args.InputPath)
		if err != nil {
			panic(err)
		}
		format, err = detectFormat(reader.(*os.File))
		if err != nil {
			panic(err)
		}
	} else {
		reader = bufio.NewReader(os.Stdin)
		format = marc.MARCXML
	}
	log.Printf("Building a new marcstream format %v", format)
	ms, err := marcdex.NewMarcStream(reader, 1000, format)
	if err != nil {
		panic(err)
	}

	indexer := marcdex.MustNewMarcIndexer(ms, recordStore)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	indexer.BatchWrite()

}
