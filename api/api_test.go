package api_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/writer"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/api"
)

func TestErrorFunc(t *testing.T) {
	w := httptest.NewRecorder()
	api.ErrorFunc(w, fmt.Errorf("err"))
	assert.Equal(t, w.Code, http.StatusBadRequest, "Status code should be 400")
	response := `{"error" : "err"}`
	assert.Equal(t, response, w.Body.String(), "Body should be json error")
}

func TestWriteResult(t *testing.T) {
	w := httptest.NewRecorder()
	result := "test"
	api.WriteResult(w, result)
	response_expected := "{\"result\":\"test\"}\n"
	assert.Equal(t, response_expected, w.Body.String(), "Body should equal")
}

func TestNewSubjectQuery(t *testing.T) {
	qs, _ := url.ParseQuery("subject=blah")
	sq := api.NewSubjectQuery(qs)
	assert.Equal(t, "blah", sq.Subject)
}

func TestQueryGraph(t *testing.T) {
	qs, _ := graph.NewQuadStore("memstore", "", nil)
	wq, err := writer.NewSingleReplication(qs, nil)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	data := quad.Make("Alice", "follows", "bob", nil)
	wq.AddQuad(data)
	w := httptest.NewRecorder()
	api.QueryGraph(qs, api.SubjectQuery{Subject: "Alice"}, w)
	assert.Equal(t, "{\"result\":[{\"object\":\"bob\",\"predicate\":\"follows\"}]}\n", w.Body.String())
}
