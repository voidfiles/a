package main

import (
	"time"

	"github.com/coreos/bbolt"
	"github.com/voidfiles/a/cli"
	"github.com/voidfiles/a/recordstore"
	"github.com/voidfiles/a/search"
)

func main() {
	args := cli.GetArgs()

	db, err := bolt.Open(args.Dbpath, 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	defer db.Close()

	if err != nil {
		panic(err)
	}
	recordStore := recordstore.MustNewRecordStore(db)
	recordStream, err := recordstore.NewRecordStream(recordStore, 1000)
	if err != nil {
		panic(err)
	}
	searchIndex := search.MustNewIndex(args.IndexPath)

	for recordStream.Next() {
		println("yo")
		records := recordStream.Value()
		searchIndex.BatchIndex(records)
	}

}
