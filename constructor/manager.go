package constructor

import (
	"encoding/json"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/example/trackerLogic/yunexpress"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/utils"
)

// [[api5, api5], [api2]]

// 存放构造器链表
type constructorStatusMap struct {
	linkList *TrackerLogicNode		// 构造器链表
	trackData *httpobj.TrackData	// 构造器需要的trackData
}

// Manager 构造器管理器
type Manager struct {
	collectorConstructorChan chan *httpobj.TrackData	// 收集器到构造器的管道
	engineConstructorChan chan *httpobj.ItemStatus		// 引擎返回的线路反馈
	ConstructorMap map[string]*constructorStatusMap		// 构造器管理器map
}

// SubmitTrackData 向构造器管理器提交 trackData
func (manager *Manager) SubmitTrackData(trackData *httpobj.TrackData) {
	manager.collectorConstructorChan <- trackData
}

// SubmitItem 向构造器管理器提交 item
func (manager *Manager) SubmitItem(item *httpobj.ItemStatus) {
	manager.engineConstructorChan <- item
}

// constructorLinkList 根据配置构造对应的构造器链表
func (manager *Manager) constructorLinkList(engine engine.Engine, trackData *httpobj.TrackData) {
	// =====================================================
	// 注册对应线路
	constructorMap := map[string]TrackerLogic{
		"api5": yunexpress.NewYunExpressLogic("API5"),
		"api2": yunexpress.NewYunExpressLogic("API2"),
	}
	// 解析构造器配置结构
	constructorList := [][]string{
		{"api5"},
		//{"api2"},
	}
	//testBaseTrackDataId := "123456789"
	// =====================================================

	// 1. 根据id, 线路列表构建线路链表
	constructorLinkList := NewTrackerLogicNode()
	for _, nodeList := range constructorList{
		var cList []TrackerLogic
		for _, node := range nodeList{
			cList = append(cList, constructorMap[node])
		}
		constructorLinkList.Push(cList)
	}

	// 构造器状态管理器
	csm := &constructorStatusMap{}

	// 单号 md5列表 作为RequestId
	trackDataStr, _ := json.Marshal(trackData)
	trackDataMd5 := utils.Md5V(string(trackDataStr))

	// 将指针移动到第一项
	constructorLinkList = constructorLinkList.Next
	for _, constructorLogic := range constructorLinkList.Logic{
		// todo: 向 redis 报告

		reqMeta := &httpobj.RequestMeta{
			RequestId: trackDataMd5,		// RequestId
			HasNextConstructorNode: constructorLinkList.HasNextNode(), // 是否存在下一个节点
		}
		req := constructorLogic.ConstructFirstRequest(reqMeta, trackData)
		// 提交第一个请求至引擎
		engine.SubmitRequests(req)
	}

	csm.linkList = constructorLinkList
	csm.trackData = trackData
	manager.ConstructorMap[trackDataMd5] = csm
}

func (manager *Manager) processItemStatus(engine engine.Engine, itemStatus *httpobj.ItemStatus) {
	// 线路正常响应
	if !itemStatus.IsError{
		// 删除状态信息
		delete(manager.ConstructorMap, itemStatus.RequestId)
		//fmt.Println("线路正常响应结束")
		// 线路异常
	}else{
		csm, ok := manager.ConstructorMap[itemStatus.RequestId]
		if !ok{
			//fmt.Println("constructorLinkList 获取失败")
			return
		}
		csmLinkList := csm.linkList

		// 不存在其他节点
		if !csmLinkList.HasNextNode(){
			delete(manager.ConstructorMap, itemStatus.RequestId)
			//fmt.Println("线路异常响应结束")
		}else{
			// 将指针移动到第一项
			csmLinkList = csmLinkList.Next
			for _, constructorLogic := range csmLinkList.Logic{
				// todo: 向 redis 报告

				reqMeta := &httpobj.RequestMeta{
					RequestId: itemStatus.RequestId,
					HasNextConstructorNode: csmLinkList.HasNextNode(),
				}
				req := constructorLogic.ConstructFirstRequest(reqMeta, csm.trackData)
				engine.SubmitRequests(req)
				//fmt.Println("线路切换")
			}
			csm.linkList = csmLinkList
		}
	}
}

func (manager *Manager) Run(engine engine.Engine) {
	// 2. 进入第一个node, 并报告redis
	// 3. 等待结果, 报告redis, 经过的线路
	for {
		select {
		// 获取到从 dataCollector 过来的单号数据
		case trackData := <- manager.collectorConstructorChan:
			// 根据单号构建对应线路链表
			manager.constructorLinkList(engine, trackData)

		// 获取从引擎回来的响应
		case itemStatus := <- manager.engineConstructorChan:
			// 处理从引擎回来的响应
			manager.processItemStatus(engine, itemStatus)
		}
	}
}

// NewConstructorManager 返回中间件管理器
func NewConstructorManager(ecchan chan *httpobj.ItemStatus) *Manager {
	return &Manager{
		collectorConstructorChan: make(chan *httpobj.TrackData),
		engineConstructorChan: ecchan,
		ConstructorMap: make(map[string]*constructorStatusMap),
	}
}