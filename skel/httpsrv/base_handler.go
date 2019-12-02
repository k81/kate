package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/k81/govalidator"
	"github.com/k81/kate"
	"github.com/k81/kate/utils"
	"go.uber.org/zap"
)

const (
	// HeaderContentLength the header name of `Content-Length`
	HeaderContentLength = "Content-Length"
	// HeaderContentType the header name of `Content-Type`
	HeaderContentType = "Content-Type"
	// MIMEApplicationJSON the application type for json
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 the application type for json of utf-8 encoding
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=UTF-8"
)

// ErrServerInternal indicates the server internal error
var ErrServerInternal = NewError(-1, "server internal error")

// BaseHandler is the enhanced version of ngs.BaseController
type BaseHandler struct{}

// ParseRequest parses and validates the api request
// nolint:lll,gocyclo
func (h *BaseHandler) ParseRequest(ctx context.Context, r *kate.Request, req interface{}) error {
	// decode json
	if r.ContentLength != 0 {
		if err := h.parseBody(req, r); err != nil {
			logger.Error("decode request", zap.Error(err))
			return err
		}
	}

	// decode query
	queryValues := r.URL.Query()
	if len(queryValues) > 0 {
		data := make(map[string]interface{})
		for key := range queryValues {
			data[key] = queryValues.Get(key)
		}

		if err := utils.Bind(req, "query", data); err != nil {
			logger.Error("bind query var failed", zap.Error(err))
			return err
		}
	}

	// decode rest var
	if len(r.RestVars) > 0 {
		data := make(map[string]interface{})
		for i := range r.RestVars {
			data[r.RestVars[i].Key] = r.RestVars[i].Value
		}

		if err := utils.Bind(req, "rest", data); err != nil {
			logger.Error("bind rest var failed", zap.Error(err))
			return err
		}
	}

	// set defaults
	if err := utils.SetDefaults(req); err != nil {
		logger.Error("set default failed", zap.Error(err))
		return ErrServerInternal
	}
	// validate
	if err := govalidator.ValidateStruct(req); err != nil {
		logger.Error("validate request", zap.Error(err))
		return err
	}
	return nil
}

// Error writes out an error response
func (h *BaseHandler) Error(ctx context.Context, w http.ResponseWriter, err interface{}) {
	Error(ctx, w, err)
}

// OK writes out a success response without data, used typically in an `update` api.
func (h *BaseHandler) OK(ctx context.Context, w http.ResponseWriter) {
	OK(ctx, w)
}

// OKData writes out a success response with data, used typically in an `get` api.
func (h *BaseHandler) OKData(ctx context.Context, w http.ResponseWriter, data interface{}) {
	OKData(ctx, w, data)
}

// EncodeJSON is a wrapper of json.Marshal()
func (h *BaseHandler) EncodeJSON(v interface{}) ([]byte, error) {
	return EncodeJSON(v)
}

// WriteJSON writes out an object which is serialized as json.
func (h *BaseHandler) WriteJSON(w http.ResponseWriter, v interface{}) error {
	return WriteJSON(w, v)
}

// parseBody 从http request 中解出json body，必须是 application/json
func (h *BaseHandler) parseBody(ptr interface{}, req *kate.Request) (err error) {
	ctype := req.Header.Get(HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		if err = utils.ParseJSON(req.RawBody, ptr); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return fmt.Errorf("unmarshal type error: expected=%v, got=%v, offset=%v",
					ute.Type, ute.Value, ute.Offset)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return fmt.Errorf("syntax error: offset=%v, error=%v",
					se.Offset, se.Error())
			} else {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported media type")
	}
	return nil
}
