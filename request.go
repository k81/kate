package kate

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/k81/kate/utils"
)

type Request struct {
	*http.Request

	Path    httprouter.Params
	RawBody []byte
}

func (r *Request) GetBool(key string) bool {
	return utils.GetBool(r.Form.Get(key))
}

func (r *Request) GetString(key string) string {
	return utils.GetString(r.Form.Get(key))
}

func (r *Request) GetInt(key string) int {
	return utils.GetInt(r.Form.Get(key))
}

func (r *Request) GetInt8(key string) int8 {
	return utils.GetInt8(r.Form.Get(key))
}

func (r *Request) GetInt16(key string) int16 {
	return utils.GetInt16(r.Form.Get(key))
}

func (r *Request) GetInt32(key string) int32 {
	return utils.GetInt32(r.Form.Get(key))
}

func (r *Request) GetInt64(key string) int64 {
	return utils.GetInt64(r.Form.Get(key))
}

func (r *Request) GetUint(key string) uint {
	return utils.GetUint(r.Form.Get(key))
}

func (r *Request) GetUint8(key string) uint8 {
	return utils.GetUint8(r.Form.Get(key))
}

func (r *Request) GetUint16(key string) uint16 {
	return utils.GetUint16(r.Form.Get(key))
}

func (r *Request) GetUint32(key string) uint32 {
	return utils.GetUint32(r.Form.Get(key))
}

func (r *Request) GetUint64(key string) uint64 {
	return utils.GetUint64(r.Form.Get(key))
}

func (r *Request) GetFloat32(key string) float32 {
	return utils.GetFloat32(r.Form.Get(key))
}

func (r *Request) GetFloat64(key string) float64 {
	return utils.GetFloat64(r.Form.Get(key))
}
