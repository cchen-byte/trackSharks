package async

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type testDataCollect struct {}
func (tc *testDataCollect) Run(collectDataChan chan *httpobj.TrackData){
	var testTrackData httpobj.TrackData
	testBaseTrackData := &httpobj.BaseTrackData{
		TrackNumber: "YT2121621236013143",
		UserId: "451",
		Lang: "en",
	}
	testTrackData = append(testTrackData, testBaseTrackData)
	collectDataChan <- &testTrackData
}
