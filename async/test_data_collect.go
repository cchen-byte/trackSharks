package async

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/httpobj"

)

type testDataCollect struct {}
func (tc *testDataCollect) Run(engine engine.Engine, trackerLogic constructor.TrackerLogic){
	var testTrackData httpobj.TrackData
	testBaseTrackData := &httpobj.BaseTrackData{
		TrackNumber: "YT2121621236013143",
		UserId: "451",
		Lang: "en",
	}
	testTrackData = append(testTrackData, testBaseTrackData)
	for i:=0; i<1; i++{
		req := trackerLogic.ConstructFirstRequest(testTrackData)
		engine.SubmitRequests(req)
	}
}
