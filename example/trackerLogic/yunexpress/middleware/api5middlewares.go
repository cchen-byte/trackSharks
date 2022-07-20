package middleware

import "github.com/cchen-byte/trackeSharkes/httpobj"

type Api5ExpressMiddleware struct {

}

func (middleware *Api5ExpressMiddleware) ProcessRequest(request *httpobj.Request) error {
	return nil
}


func (middleware *Api5ExpressMiddleware) ProcessResponse(response *httpobj.Response) error {

	return nil
}
