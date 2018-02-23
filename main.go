package main

import (
	_ "github.com/cayleygraph/cayley/graph/kv/bolt" // Makes sure bolt is included
)

func main() {
	// args := cli.GetArgs()
	// qs, err := graph.NewQuadStore(args.Db, args.Dbpath, graph.Options{"nosync": args.Nosync})
	// if err != nil {
	// 	panic(err)
	// }
	// index := search.MustNewIndex(args.IndexPath)
	//
	// resolver := authority.NewResolver(qs, index)
	//
	// log.Printf("Starting up an http server")
	// mux := api.NewApi(resolver)
	//
	// address := args.IP + ":" + args.Port
	// log.Fatal(http.ListenAndServe(address, mux))

}
