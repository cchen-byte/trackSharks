package constructor

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

// TrackerLogic 追踪逻辑, 由callback链接
type TrackerLogic interface {
	ConstructFirstRequest(httpobj.TrackData) *httpobj.Request 	// 根据数据构造请求
}

// TrackerLogicNode 追踪逻辑节点, 一个节点为一个线路
type TrackerLogicNode struct {
	Logic TrackerLogic
	Next *TrackerLogicNode
}