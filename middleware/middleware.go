package middleware

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type Middlewares interface {
	ProcessRequest(request *httpobj.Request) error
	ProcessResponse(response *httpobj.Response) error
}