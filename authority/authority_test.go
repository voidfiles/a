package authority_test

import (
	"log"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/writer"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/authority"
	"github.com/voidfiles/a/search"
)

type AuthoritySearchMock struct{}

func (a *AuthoritySearchMock) BatchIndex([]search.QuadForIndex) error {
	return nil
}

func (a *AuthoritySearchMock) Query(label string) (*bleve.SearchResult, error) {
	return nil, nil
}

func TestNewQueryExecutor(t *testing.T) {
	qs, _ := graph.NewQuadStore("memstore", "", graph.Options{})
	index := &AuthoritySearchMock{}
	resolver := authority.NewResolver(qs, index)
	w, err := writer.NewSingleReplication(qs, graph.Options{})
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	data := quad.Make("alice", "follows", "bob", "")
	w.AddQuad(data)
	q := "alice"
	output, err := resolver.FindLabelsForID(q)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	expected := []authority.PredicateObject(
		[]authority.PredicateObject{
			authority.PredicateObject{
				Predicate: "follows",
				Object:    "bob",
			},
		},
	)
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

	assert.Equal(t, query_expected, authority.BuildLabelsForIdQuery("abc"), "Body should equal")
}
