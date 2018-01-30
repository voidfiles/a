package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/query"
)

// ErrorFunc handles writing an error response
func ErrorFunc(w query.ResponseWriter, err error) {
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
			ErrorFunc(w, err)
			return
		}

		session.Collate(res)
	}
	output, err := session.Results()

	if err != nil {
		ErrorFunc(w, err)
		return
	}
	_ = WriteResult(w, output)
}

type SubjectQuery struct {
	Subject string
}

func NewSubjectQuery(req *http.Request) SubjectQuery {
	return SubjectQuery{
		Subject: req.URL.Query().Get("subject"),
	}
}

func SubjectQueryHandler(qs graph.QuadStore, w http.ResponseWriter, req *http.Request) {
	graphQuery := NewSubjectQuery(req)
	log.Printf("Looking for blah %s", graphQuery)
	QueryGraph(qs, graphQuery, w)
}

func wrap(qs graph.QuadStore, handler func(graph.QuadStore, http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(qs, w, req)
		return
	}
}

func NewApi(qs graph.QuadStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/query/subject", wrap(qs, SubjectQueryHandler))

	return mux
}
