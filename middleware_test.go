package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareAllMethods(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPut, http.MethodPatch:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte("response"))
	})

	wrapped := Middleware(handler)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/"+strings.ToLower(method), nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)
	}

	out := buf.String()
	t.Log("\n" + out)

	for _, method := range methods {
		if !strings.Contains(out, method) {
			t.Errorf("expected %s in output", method)
		}
	}
}

func TestResponseWriterWrapper(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapper := &responseWriterWrapper{
		ResponseWriter: rec,
		status:         http.StatusOK,
	}

	wrapper.WriteHeader(http.StatusNotFound)
	wrapper.Write([]byte("not found"))
	wrapper.WriteHeader(http.StatusOK)

	if wrapper.status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", wrapper.status)
	}
	if wrapper.size != 9 {
		t.Errorf("expected size 9, got %d", wrapper.size)
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected underlying status 404, got %d", rec.Code)
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int
		expected string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0K"},
		{1536, "1.5K"},
		{1048576, "1.0M"},
		{1572864, "1.5M"},
	}

	for _, tt := range tests {
		got := formatSize(tt.bytes)
		if got != tt.expected {
			t.Errorf("formatSize(%d) = %s, want %s", tt.bytes, got, tt.expected)
		}
	}
}

func TestMiddlewareLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelWarn)
	defer SetLevel(LevelDebug)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Middleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if buf.Len() > 0 {
		t.Error("middleware logs should be filtered at WARN level")
	}
}
