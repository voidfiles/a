package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/cayleygraph/cayley/query"
	"github.com/voidfiles/a/authority"
)

// AuthorityResolver provides the interface of authority.Resolver
type AuthorityResolver interface {
	FindLabelsForID(string) ([]authority.PredicateObject, error)
	FindIdsFromLabel(string) (*bleve.SearchResult, error)
}

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

// SubjectQueryHandler provides an HTTP Api to resolver.FindLabelsForID
func SubjectQueryHandler(resolver AuthorityResolver, w http.ResponseWriter, req *http.Request) {
	resp, err := resolver.FindLabelsForID(req.URL.Query().Get("subject"))
	if err != nil {
		ErrorFunc(w, err)
	}
	_ = WriteResult(w, resp)
}

// ObjectQueryHandler provides an HTTP Api to resolver.FindLabelsForID
func ObjectQueryHandler(resolver AuthorityResolver, w http.ResponseWriter, req *http.Request) {
	resp, err := resolver.FindIdsFromLabel(req.URL.Query().Get("object"))
	if err != nil {
		ErrorFunc(w, err)
	}
	_ = WriteResult(w, resp)
}

func wrap(resolver AuthorityResolver, handler func(AuthorityResolver, http.ResponseWriter, *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(resolver, w, req)
		return
	}
}

// NewApi creates an http handler
func NewApi(resolver AuthorityResolver) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/query/subject", wrap(resolver, SubjectQueryHandler))
	mux.HandleFunc("/api/v1/query/object", wrap(resolver, ObjectQueryHandler))
	return mux
}
