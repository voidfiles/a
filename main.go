package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	_ "github.com/cayleygraph/cayley/graph/sql/postgres"
	_ "github.com/cayleygraph/cayley/query/gizmo"
	"github.com/voidfiles/a/api"
)

type cliOptions struct {
	db     string
	dbpath string
	nosync *bool
	ip     string
	port   string
}

func getEnvString(key string, fallback *string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = *fallback
	}
	return value
}

func getArgs() cliOptions {
	db := flag.String("db", "bolt", "Database Backend. (default \"bolt\")")
	dbpath := flag.String("dbpath", "/tmp/testdb", "Path to the database. (default \"/tmp/testdb\")")
	nosync := flag.Bool("nosync", true, "Should db not sync to disk (default true)")
	ip := flag.String("ip", "127.0.0.1", "IP server should bind too (default 127.0.0.1)")
	port := flag.String("port", "8080", "Port server should bind too (default 8080)")

	flag.Parse()

	return cliOptions{
		db:     getEnvString("A_DB", db),
		dbpath: getEnvString("A_DBPATH", dbpath),
		nosync: nosync,
		ip:     getEnvString("A_IP", ip),
		port:   getEnvString("PORT", port),
	}
}

// NewStore constructs a new quadstore based on cli args
func NewStore(args cliOptions) (graph.QuadStore, error) {
	log.Printf("Connecting to db %s dbpath %s", args.db, args.dbpath)
	return graph.NewQuadStore(
		args.db,
		args.dbpath,
		graph.Options{"nosync": args.nosync},
	)
}

func main() {
	args := getArgs()
	qs, err := NewStore(args)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting up an http server")
	mux := api.NewApi(qs)

	address := args.ip + ":" + args.port
	log.Fatal(http.ListenAndServe(address, mux))

}
