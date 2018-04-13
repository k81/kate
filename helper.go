package kate

import (
	"context"

	"github.com/k81/kate/log"
)

func Error(ctx context.Context, w ResponseWriter, errInfo ErrorInfo) {
	result := &Result{
		Status: errInfo.Code(),
		Msg:    errInfo.Error(),
	}

	if errInfoWithData, ok := errInfo.(ErrorInfoWithData); ok {
		result.Data = errInfoWithData.Data()
	}

	if err := w.WriteJSON(result); err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}

func Ok(ctx context.Context, w ResponseWriter) {
	OkData(ctx, w, nil)
}

func OkData(ctx context.Context, w ResponseWriter, data interface{}) {
	result := &Result{
		Status: ErrSuccess.Code(),
		Msg:    ErrSuccess.Error(),
		Data:   data,
	}

	if err := w.WriteJSON(result); err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}

func WriteXML(ctx context.Context, w ResponseWriter, data interface{}) {
	if err := w.WriteXML(data); err != nil {
		log.Error(ctx, "write xml response", "error", err)
	}
}

func WriteJSON(ctx context.Context, w ResponseWriter, data interface{}) {
	if err := w.WriteJSON(data); err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}
