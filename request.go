package kate

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request defines the http request
type Request struct {
	*http.Request

	RestVars httprouter.Params
	RawBody  []byte
}
