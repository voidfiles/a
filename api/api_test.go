package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
