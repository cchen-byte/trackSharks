package constructor

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

// TrackerLogic 追踪逻辑, 由callback链接
type TrackerLogic interface {
	ConstructFirstRequest(*httpobj.RequestMeta, *httpobj.TrackData) *httpobj.Request 	// 根据数据构造请求
}

// TrackerLogicNode 追踪逻辑节点, 一个节点为一个线路
type TrackerLogicNode struct {
	Logic []TrackerLogic
	Next *TrackerLogicNode
}

// Push 尾部插入新的节点
func (tln *TrackerLogicNode) Push(trackLogicList []TrackerLogic) {
	// 找到最后一个节点, 再插入
	func(node *TrackerLogicNode, trackLogicList []TrackerLogic) {
		for {
			if node.Next == nil{
				node.Next =  &TrackerLogicNode{
						Logic: trackLogicList,
						Next: nil}
				break
			}
			node = node.Next
		}
	}(tln, trackLogicList)
}

// HasNextNode 是否还存在下一个节点
func (tln *TrackerLogicNode) HasNextNode() bool {
	if tln.Next != nil{
		return true
	}
	return false
}

// NewTrackerLogicNode 返回头节点
func NewTrackerLogicNode() *TrackerLogicNode {
	return &TrackerLogicNode{}
}