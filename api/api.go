package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

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
	executor := NewQueryExecutor(qs)
	output, err := executor.Execute(q)
	if err != nil {
		ErrorFunc(w, err)
	}

	_ = WriteResult(w, output)
}

type SubjectQuery struct {
	Subject string
}

func NewSubjectQuery(qs url.Values) SubjectQuery {
	return SubjectQuery{
		Subject: qs.Get("subject"),
	}
}

func SubjectQueryHandler(qs graph.QuadStore, w http.ResponseWriter, req *http.Request) {
	graphQuery := NewSubjectQuery(req.URL.Query())
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
