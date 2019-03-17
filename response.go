package kate

import (
	"net/http"
)

// ResponseWriter defines the response writer
type ResponseWriter interface {
	http.ResponseWriter

	StatusCode() int

	RawBody() []byte
}

type responseWriter struct {
	http.ResponseWriter

	wroteHeader bool
	statusCode  int
	rawBody     []byte
}

func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

func (w *responseWriter) RawBody() []byte {
	return w.rawBody
}

func (w *responseWriter) WriteHeader(code int) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
	w.statusCode = code
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.rawBody = b
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	// nolint:errcheck
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}
