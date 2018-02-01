package api_test

import (
	"log"
	"testing"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
	_ "github.com/cayleygraph/cayley/query/gizmo"
	"github.com/cayleygraph/cayley/writer"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/api"
)

func TestNewQueryExecutor(t *testing.T) {
	qs, _ := graph.NewQuadStore("memstore", "", nil)
	api_exec := api.NewQueryExecutor(qs)
	w, err := writer.NewSingleReplication(qs, nil)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	data := quad.Make("alice", "follows", "bob", nil)
	w.AddQuad(data)
	q := api.SubjectQuery{Subject: "alice"}
	output, err := api_exec.Execute(q)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	expected := []interface{}{
		map[string]interface{}{
			"predicate": "follows",
			"object":    "bob",
		},
	}
	assert.Equal(t, expected, output)
}

func TestBuildQuery(t *testing.T) {
	query_expected := `var vID = "abc"; g.V(vID).OutPredicates().ForEach( function(r){
		g.V(vID).Out(r.id).ForEach( function(t){
			var node = {
			  predicate: r.id,
			  object: t.id
			}
			g.Emit(node)
		})
	});`
	q := api.SubjectQuery{
		Subject: "abc",
	}
	assert.Equal(t, query_expected, api.BuildQuery(q), "Body should equal")
}
