package main

import (
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt" // Makes sure bolt is included
	"github.com/voidfiles/a/authority"
	"github.com/voidfiles/a/cli"
	"github.com/voidfiles/a/search"
)

func main() {
	args := cli.GetArgs()
	qs, err := graph.NewQuadStore(args.Db, args.Dbpath, graph.Options{"nosync": args.Nosync})
	if err != nil {
		panic(err)
	}
	index := search.MustNewIndex(args.IndexPath)

	indexer := authority.MustNewIndexer(qs, index)
	indexer.ProcessFullText()
}
