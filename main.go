package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	_ "github.com/cayleygraph/cayley/query/gizmo"
	"github.com/voidfiles/a/api"
)

type cliOptions struct {
	db *string
}

func getArgs() cliOptions {
	db := flag.String("db", "", "Path to the database")

	flag.Parse()
	return cliOptions{db: db}
}

func NewStore(args cliOptions) (graph.QuadStore, error) {
	log.Printf("Connecting to local db %s", *args.db)
	return graph.NewQuadStore("bolt", *args.db, graph.Options{"nosync": true})
}

func main() {
	args := getArgs()
	qs, err := NewStore(args)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting up an http server")
	mux := api.NewApi(qs)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))

}
