package kate

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	StatusCode() int

	RawBody() []byte

	EncodeJSON(v interface{}) ([]byte, error)

	WriteJSON(v interface{}) error

	EncodeXML(v interface{}) ([]byte, error)

	WriteXML(v interface{}) error
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

func (w *responseWriter) EncodeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (w *responseWriter) WriteJSON(v interface{}) error {
	b, err := w.EncodeJSON(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *responseWriter) EncodeXML(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (w *responseWriter) WriteXML(v interface{}) error {
	b, err := w.EncodeXML(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
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
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}
