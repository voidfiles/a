package api

import (
	"context"
	"fmt"
	"log"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/query"
)

// NewQueryExecutor creates and returns a interface
func NewQueryExecutor(qs graph.QuadStore) QueryExecutor {
	l := query.GetLanguage("gizmo")
	return QueryExecutor{
		Session:       l.HTTP(qs),
		ResultChannel: make(chan query.Result, 5),
	}
}

// BuildQuery creates a gizmo query based on a subject
func BuildQuery(q SubjectQuery) string {
	quStart := `g.V(vID).OutPredicates().ForEach( function(r){
		g.V(vID).Out(r.id).ForEach( function(t){
			var node = {
			  predicate: r.id,
			  object: t.id
			}
			g.Emit(node)
		})
	});`

	return fmt.Sprintf("var vID = \"%s\"; %s", q.Subject, quStart)
}

// QueryExecutor will execute a gizmo query against a QuadStore
type QueryExecutor struct {
	Session       query.HTTP
	ResultChannel chan query.Result
}

// Execute will execute a query against a quad store
func (qe *QueryExecutor) Execute(q SubjectQuery) (interface{}, error) {
	ctx := context.TODO()
	qu := BuildQuery(q)
	log.Printf("Executing query %s", qu)
	go qe.Session.Execute(ctx, qu, qe.ResultChannel, 100)
	for res := range qe.ResultChannel {
		if err := res.Err(); err != nil {
			if err == nil {
				continue // wait for results channel to close
			}

			return nil, err
		}

		qe.Session.Collate(res)
	}
	output, err := qe.Session.Results()
	if err != nil {
		return nil, err
	}

	return output, nil
}
