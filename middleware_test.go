package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	myHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	tests := []struct {
		description string
		key         string
		contType    string
		expCode     int
		expRes      string
	}{
		{
			"success case", "a601e44e306e430f8dde987f65844f05", "application/json",
			200, "hello world",
		},
		{
			"Error case: incorrect authentication key", "a601e44e306e430f8dde987f65844000",
			"application/json", 401, "authentication failed",
		},
		{
			"Error case: incorrect content type", "a601e44e306e430f8dde987f65844f05",
			"", 415, "Header Content-Type incorrect",
		},
	}
	for i, tc := range tests {
		req := httptest.NewRequest("GET", "/hello", nil)
		req.Header.Set("X-API-KEY", tc.key)
		req.Header.Set("Content-Type", tc.contType)

		resRec := httptest.NewRecorder()

		middleware(myHandlerFunc).ServeHTTP(resRec, req)

		assert.Equal(t, tc.expCode, resRec.Code, "Test[%d] failed\n(%s)", i, tc.description)
		assert.Equal(t, tc.expRes, resRec.Body.String(), "Test[%d] failed\n(%s)", i, tc.description)
	}
}
