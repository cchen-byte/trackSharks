package downloader

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type Downloader interface {
	Fetch(request *httpobj.Request) (*httpobj.Response, error)
}