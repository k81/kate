package kate

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Request struct {
	*http.Request

	RestVars httprouter.Params
	RawBody  []byte
}
