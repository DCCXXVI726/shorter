package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShorterHandler(t *testing.T) {
	tests := []struct {
		name       string
		handler    func(http.ResponseWriter, *http.Request)
		url        string
		statusCode int
	}{
		{
			name:       "Simple shorter test",
			handler:    shorterHandler,
			url:        "http://127.0.0.1:8080/a/?url=http%3A%2F%2Fgoogle.com%2F%3Fq%3Dgolang",
			statusCode: http.StatusOK,
		},
		{
			name:       "shorter empty url",
			handler:    shorterHandler,
			url:        "http://127.0.0.1:8080/a/?url=",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Simple redirect test",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/bZ2EjaV1",
			statusCode: http.StatusFound,
		},
		{
			name:       "short token",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/short",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "long token",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/longggggg",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "invalid token",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/kekwait$",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "not found",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/kekwaitt",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "empty token",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "to many parameters",
			handler:    redirectHandler,
			url:        "http://127.0.0.1:8080/s/to/many",
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.url, nil)
			w := httptest.NewRecorder()
			tc.handler(w, req)
			result := w.Result()
			if result.StatusCode != tc.statusCode {
				t.Errorf("handlerShorter StatusCode = %d, want %d", result.StatusCode, tc.statusCode)
			}
		})
	}
}
