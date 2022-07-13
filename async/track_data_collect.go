package async

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type TrackDataCollect interface {
	Run(engine engine.Engine, trackerLogic constructor.TrackerLogic)
}

type testDataCollect struct {}
func (tc *testDataCollect) Run(engine engine.Engine, trackerLogic constructor.TrackerLogic){
	testTrackData := httpobj.TrackData{}
	for i:=0; i<1; i++{
		req := trackerLogic.ConstructFirstRequest(testTrackData)
		engine.SubmitRequests(req)
	}
}