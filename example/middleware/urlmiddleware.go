package middleware

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type HBUrlMiddleware struct {
}

func (middleware *HBUrlMiddleware) ProcessRequest(request *httpobj.Request) error {
	request.Url = "https://httpbin.org/get"
	return nil
}

func (middleware *HBUrlMiddleware) ProcessResponse(response *httpobj.Response) error {
	return nil
}
