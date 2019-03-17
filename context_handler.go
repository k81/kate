package kate

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/k81/log"
)

// ContextHandler defines the handler interface
type ContextHandler interface {
	ServeHTTP(context.Context, ResponseWriter, *Request)
}

// ContextHandlerFunc defines the handler func adapter
type ContextHandlerFunc func(context.Context, ResponseWriter, *Request)

// ServeHTTP implements the ContextHandler interface
func (h ContextHandlerFunc) ServeHTTP(ctx context.Context, w ResponseWriter, r *Request) {
	h(ctx, w, r)
}

// Handle adapte the ContextHandler to httprouter.Handle func
func Handle(ctx context.Context, h ContextHandler, maxBodyBytes int64) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var (
			request        *Request
			response       *responseWriter
			err            error
			newctx, cancel = context.WithCancel(ctx)
		)

		defer cancel()

		request = &Request{
			Request:  r,
			RestVars: params,
		}

		response = &responseWriter{
			ResponseWriter: w,
			wroteHeader:    false,
		}

		if maxBodyBytes > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		}

		if request.RawBody, err = ioutil.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// nolint:errcheck
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		// nolint:errcheck
		r.Body.Close()

		r.Body = ioutil.NopCloser(bytes.NewReader(request.RawBody))

		err = r.ParseMultipartForm(maxBodyBytes)
		switch {
		case err == http.ErrNotMultipart:
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			// nolint:errcheck
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			log.Error(ctx, "read request", "error", err)
			return
		}

		h.ServeHTTP(newctx, response, request)
	}
	return httprouter.Handle(f)
}

// StdHandler adapte ContextHandler to http.Handler interface
func StdHandler(ctx context.Context, h ContextHandler, maxBodyBytes int64) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var (
			request        *Request
			response       *responseWriter
			err            error
			newctx, cancel = context.WithCancel(ctx)
		)

		defer cancel()

		request = &Request{
			Request: r,
		}

		response = &responseWriter{
			ResponseWriter: w,
			wroteHeader:    false,
		}

		if maxBodyBytes > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		}

		if request.RawBody, err = ioutil.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// nolint:errcheck
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		// nolint:errcheck
		r.Body.Close()

		r.Body = ioutil.NopCloser(bytes.NewReader(request.RawBody))

		err = r.ParseMultipartForm(maxBodyBytes)
		switch {
		case err == http.ErrNotMultipart:
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			// nolint:errcheck
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			log.Error(ctx, "read request", "error", err)
			return
		}

		h.ServeHTTP(newctx, response, request)
	}
	return http.HandlerFunc(f)
}
