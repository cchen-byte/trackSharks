package async

import "github.com/cchen-byte/trackeSharkes/httpobj"

type TrackDataCollect interface {
	Run(chan *httpobj.TrackData)
}
