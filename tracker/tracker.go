package tracker

import "github.com/cchen-byte/trackeSharkes/httpobj"

type Tracker interface {
	fetchWork(request *httpobj.Request) (*httpobj.Response, error)
	parseWork(request *httpobj.Request, resp *httpobj.Response) (*httpobj.ParseResult, error)
	Run()
}
