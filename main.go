package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/query"
	_ "github.com/cayleygraph/cayley/query/gizmo"
)

type cliOptions struct {
	db *string
}

type Node struct {
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

type Response struct {
	Nodes []Node `json:"nodes"`
}

func newResponse() Response {
	return Response{}
}

func (r *Response) addNode(predicate, object string) {
	node := Node{
		Predicate: predicate,
		Object:    object,
	}
	r.Nodes = append(r.Nodes, node)
}

func getArgs() cliOptions {
	db := flag.String("db", "", "Path to the database")

	flag.Parse()
	return cliOptions{db: db}
}

func jsonResponse(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(value)
	w.Write(data)
}

func NewStore(args cliOptions) (graph.QuadStore, error) {
	log.Printf("Connecting to local db %s", *args.db)
	return graph.NewQuadStore("bolt", *args.db, graph.Options{"nosync": true})
}

type SubjectQuery struct {
	Subject string
}

func NewSubjectQuery(req *http.Request) SubjectQuery {
	return SubjectQuery{
		Subject: req.URL.Query().Get("subject"),
	}
}

func defaultErrorFunc(w query.ResponseWriter, err error) {
	data, _ := json.Marshal(err.Error())
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error" : `))
	w.Write(data)
	w.Write([]byte(`}`))
}

type SuccessQueryWrapper struct {
	Result interface{} `json:"result"`
}

func WriteResult(w io.Writer, result interface{}) error {
	enc := json.NewEncoder(w)
	//enc.SetIndent("", " ")
	return enc.Encode(SuccessQueryWrapper{result})
}

func QueryGraph(qs graph.QuadStore, q SubjectQuery, w http.ResponseWriter) {
	ctx := context.TODO()
	l := query.GetLanguage("gizmo")
	session := l.HTTP(qs)
	c := make(chan query.Result, 5)
	quStart := `g.V(vID).OutPredicates().ForEach( function(r){
		g.V(vID).Out(r.id).ForEach( function(t){
			var node = {
			  source: r.id,
			  target: t.id
			}
			g.Emit(node)
		})
	});`
	qu := fmt.Sprintf("var vID = \"%s\"; %s", q.Subject, quStart)
	log.Printf("Executing query %s", qu)
	go session.Execute(ctx, qu, c, 100)
	for res := range c {
		if err := res.Err(); err != nil {
			if err == nil {
				continue // wait for results channel to close
			}
			defaultErrorFunc(w, err)
			return
		}

		session.Collate(res)
	}
	output, err := session.Results()

	if err != nil {
		defaultErrorFunc(w, err)
		return
	}
	_ = WriteResult(w, output)
}

func main() {
	args := getArgs()
	qs, err := NewStore(args)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting up an http server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		graphQuery := NewSubjectQuery(req)
		log.Printf("Looking for blah %s", graphQuery)
		QueryGraph(qs, graphQuery, w)
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))

}
