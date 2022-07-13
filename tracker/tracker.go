package tracker

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
)

type Tracker interface {
	fetchWork(request *httpobj.Request) (*httpobj.Response, error)
	parseWork(request *httpobj.Request, resp *httpobj.Response) (*httpobj.ParseResult, error)
	Run(trackerMiddlewaresManager *middleware.Manager, request *httpobj.Request) (*httpobj.ParseResult, error)
}
