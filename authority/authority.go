package authority

import (
	"context"
	"fmt"
	"log"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/query"
	_ "github.com/cayleygraph/cayley/query/gizmo"
	"github.com/voidfiles/a/search"
)

func BuildLabelsForIdQuery(id string) string {
	quStart := `g.V(vID).OutPredicates().ForEach( function(r){
		g.V(vID).Out(r.id).ForEach( function(t){
			var node = {
			  predicate: r.id,
			  object: t.id
			}
			g.Emit(node)
		})
	});`

	return fmt.Sprintf("var vID = \"%s\"; %s", id, quStart)
}

func NewResolver(qs graph.QuadStore, index search.AuthoritySearch) *Resolver {
	return &Resolver{
		qs,
		index,
	}
}

type Resolver struct {
	qs    graph.QuadStore
	index search.AuthoritySearch
}

type Quad struct {
	Subject   string
	Predicate string
	Object    string
}

type QuadValue struct {
	Value graph.Value
	Quad  quad.Quad
}

func (r *Resolver) ExecuteGizmoQuery(gizmoQuery string) ([]interface{}, error) {
	resultChannel := make(chan query.Result, 5)

	l := query.GetLanguage("gizmo")

	session := l.Session(r.qs)

	ctx := context.TODO()

	log.Printf("Executing query %s", gizmoQuery)

	go session.Execute(ctx, gizmoQuery, resultChannel, 100)
	results := make([]interface{}, 0)
	for res := range resultChannel {
		if err := res.Err(); err != nil {
			if err == nil {
				continue // wait for results channel to close
			}

			return nil, err
		}

		results = append(results, res.Result())

	}

	return results, nil
}

type PredicateObject struct {
	Predicate string
	Object    string
}

func (r *Resolver) FindLabelsForID(id string) ([]PredicateObject, error) {
	gizmoQuery := BuildLabelsForIdQuery(id)
	output, err := r.ExecuteGizmoQuery(gizmoQuery)
	if err != nil {
		panic(err)
	}
	results := make([]PredicateObject, 0)
	log.Print(output)
	for _, item := range output {
		itemMap := item.(map[string]interface{})
		results = append(results, PredicateObject{
			itemMap["predicate"].(string),
			itemMap["object"].(string),
		})

	}

	return results, nil
}

func (r *Resolver) FindIdsFromLabel(label string) (interface{}, error) {
	return "", nil
}
