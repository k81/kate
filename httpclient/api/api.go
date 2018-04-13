package api

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/k81/kate/httpclient"
	"github.com/k81/kate/log"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"
	StatusOK        = 0
)

type ApiError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("API Error: code=%d, msg=%s", e.Code, e.Msg)
}

type ApiResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}

func TryCallJSON(ctx context.Context, url string, req interface{}, resp interface{}, timeout time.Duration, backoff []time.Duration) (err error) {
	var (
		reqVal  = reflect.ValueOf(req)
		reqBody []byte
	)

	if reqVal.IsValid() && !reqVal.IsNil() {
		if reqBody, err = json.Marshal(req); err != nil {
			log.Error(ctx, "marshal req", "error", err)
			return
		}
	}

	return TryCall(ctx, "POST", url, ContentTypeJSON, string(reqBody), resp, timeout, backoff)
}

func TryCall(ctx context.Context, method, url, contentType, reqBody string, resp interface{}, timeout time.Duration, backoff []time.Duration) (err error) {
	var (
		apiResp  = &ApiResponse{}
		respVal  = reflect.ValueOf(resp)
		respBody string
	)

	if respBody, err = httpclient.TryCall(ctx, method, url, contentType, reqBody, timeout, backoff); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(respBody), apiResp); err != nil {
		log.Error(ctx, "unmarshal api response", "error", err)
		return err
	}

	if apiResp.Code != StatusOK {
		return &ApiError{Code: apiResp.Code, Msg: apiResp.Msg}
	}

	if respVal.IsValid() && !respVal.IsNil() {
		if err = json.Unmarshal(apiResp.Data, resp); err != nil {
			log.Error(ctx, "unmarshal api response body", "content", string(apiResp.Data), "error", err)
			return err
		}
	}

	return
}
