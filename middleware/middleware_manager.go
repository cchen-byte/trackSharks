package middleware

import (
	"github.com/cchen-byte/trackeSharkes/utils"
	"sync"
)

type Manager struct {
	MiddlewaresMap sync.Map
}


// GetMiddlewares 根据request的requestMiddleware获取对应的middleware的切片
func (manager *Manager) GetMiddlewares(requestMiddleware map[string]int) ([]Middlewares, error) {
	// 请求无中间件参数
	if len(requestMiddleware) == 0 {
		return []Middlewares{}, nil
	}

	// 若已经注册,则直接返回
	if middlewareList, ok := manager.MiddlewaresMap.Load(&requestMiddleware); ok {
		return middlewareList.([]Middlewares), nil
	}

	// 排序
	requestMiddlewareMapList := utils.SortMapByValue(requestMiddleware)
	var requestMiddlewareList []Middlewares
	for _, v := range requestMiddlewareMapList{
		if middleware, ok := manager.MiddlewaresMap.Load(v.Key); ok {
			requestMiddlewareList = append(requestMiddlewareList, middleware.(Middlewares))
		}

	}
	// 注册
	manager.MiddlewaresMap.Store(&requestMiddleware, requestMiddlewareList)
	return requestMiddlewareList, nil
}

// NewMiddlewareManager 返回中间件管理器
func NewMiddlewareManager() *Manager {
	return &Manager{
		MiddlewaresMap: sync.Map{},
	}
}


var TrackerMiddlewaresManager = NewMiddlewareManager()
