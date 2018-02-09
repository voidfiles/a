package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/api"
	"github.com/voidfiles/a/authority"
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

type ResolverMock struct{}

func (r *ResolverMock) FindLabelsForID(q string) ([]authority.PredicateObject, error) {
	return []authority.PredicateObject(
		[]authority.PredicateObject{
			authority.PredicateObject{
				Predicate: "follows",
				Object:    "bob",
			},
		},
	), nil
}

func TestQueryGraph(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"https://example.com/api/v1/query/subject?subject=Alice",
		nil,
	)
	api.SubjectQueryHandler(&ResolverMock{}, w, req)
	assert.Equal(t, "{\"result\":[{\"Predicate\":\"follows\",\"Object\":\"bob\"}]}\n", w.Body.String())
}
