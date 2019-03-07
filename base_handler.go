package kate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k81/govalidator"
	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

const (
	HeaderContentLength            = "Content-Length"
	HeaderContentType              = "Content-Type"
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=UTF-8"
)

var ErrServerInternal = NewError(-1, "server internal error")

// APIResponse represents the api response body
type APIResponse struct {
	ErrNO  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

// BaseHandler is the enhanced version of ngs.BaseController
type BaseHandler struct{}

// ParseRequest parses and validates the api request
// nolint:lll,gocyclo
func (h *BaseHandler) ParseRequest(ctx context.Context, r *Request, req interface{}) error {
	// decode json
	if r.ContentLength != 0 {
		if err := h.parseBody(req, r); err != nil {
			log.Error(ctx, "decode request", "error", err)
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
			log.Error(ctx, "bind query var failed", "error", err)
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
			log.Error(ctx, "bind rest var failed", "error", err)
			return err
		}
	}

	// validate
	if err := govalidator.ValidateStruct(req); err != nil {
		log.Error(ctx, "validate request", "error", err)
		return err
	}
	return nil
}

// Error writes out an error response
func (h *BaseHandler) Error(ctx context.Context, w ResponseWriter, err interface{}) {
	errInfo, ok := err.(ErrorInfo)
	if !ok {
		errInfo = ErrServerInternal
	}

	apiResp := &APIResponse{
		ErrNO:  errInfo.Code(),
		ErrMsg: errInfo.Error(),
	}

	if errInfoWithData, ok := errInfo.(ErrorInfoWithData); ok {
		apiResp.Data = errInfoWithData.Data()
	}

	if err := h.WriteJSON(w, apiResp); err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}

// OK writes out a success response without data, used typically in an `update` api.
func (h *BaseHandler) OK(ctx context.Context, w ResponseWriter) {
	h.OKData(ctx, w, nil)
}

// OKData writes out a success response with data, used typically in an `get` api.
func (h *BaseHandler) OKData(ctx context.Context, w ResponseWriter, data interface{}) {
	apiResp := &APIResponse{
		ErrNO:  0,
		ErrMsg: "success",
		Data:   data,
	}

	if err := h.WriteJSON(w, apiResp); err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}

// EncodeJSON is a wrapper of json.Marshal()
func (h *BaseHandler) EncodeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// WriteJSON writes out an object which is serialized as json.
func (h *BaseHandler) WriteJSON(w ResponseWriter, v interface{}) error {
	b, err := h.EncodeJSON(v)
	if err != nil {
		return err
	}
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

// parseBody 从http request 中解出json body，必须是 application/json
func (h *BaseHandler) parseBody(ptr interface{}, req *Request) (err error) {
	ctype := req.Header.Get(HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		if err = utils.ParseJSON(bytes.NewReader(req.RawBody), ptr); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return fmt.Errorf("Unmarshal type error: expected=%v, got=%v, offset=%v",
					ute.Type, ute.Value, ute.Offset)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return fmt.Errorf("Syntax error: offset=%v, error=%v",
					se.Offset, se.Error())
			} else {
				return err
			}
		}
	default:
		return fmt.Errorf("Unsupported media type")
	}
	return nil
}
