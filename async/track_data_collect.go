package async

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/engine"
)

type TrackDataCollect interface {
	Run(engine engine.Engine, trackerLogic constructor.TrackerLogic)
}
