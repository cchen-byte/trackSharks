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

	testBaseTrackData2 := &httpobj.BaseTrackData{
		TrackNumber: "YT2121721236014632",
		UserId: "451",
		Lang: "en",
	}
	testTrackData = append(testTrackData, testBaseTrackData)
	testTrackData = append(testTrackData, testBaseTrackData2)
	//testTrackData = append(testTrackData, testBaseTrackData2)
	//testTrackData = append(testTrackData, testBaseTrackData2)
	//testTrackData = append(testTrackData, testBaseTrackData2)
	collectDataChan <- &testTrackData
}
