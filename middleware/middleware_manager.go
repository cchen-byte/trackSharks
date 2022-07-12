package middleware

import (
	"encoding/json"
	"github.com/cchen-byte/trackeSharkes/utils"
	"strings"
	"sync"
)

type Manager struct {
	MiddlewaresMap sync.Map
}

// AddMiddleware 添加中间件
func (manager *Manager) AddMiddleware(middlewareType string, mw Middlewares) {
	middlewaresMiddlewareType, _ := manager.MiddlewaresMap.Load(middlewareType)
	manager.MiddlewaresMap.Store(middlewareType, append(middlewaresMiddlewareType.([]Middlewares), mw))
}

// GetMiddlewares 根据request的requestMiddleware获取对应的middleware的切片
func (manager *Manager) GetMiddlewares(requestMiddleware map[string]int, trackerMiddleware map[string]Middlewares) ([]Middlewares, error) {
	// 请求无中间件参数
	if len(requestMiddleware) == 0 {
		return []Middlewares{}, nil
	}

	// 将请求需要的中间件注册到MiddlewaresMap中, key为请求中间件md5值, value为对应的中间件切片
	marshal, err := json.Marshal(requestMiddleware)
	if err != nil {
		return nil, err
	}
	mapName := utils.Md5V(string(marshal))

	// 若已经注册,则直接返回
	if middlewareList, ok := manager.MiddlewaresMap.Load(mapName); ok {
		return middlewareList.([]Middlewares), nil
	}

	// 注册
	var middlewares []Middlewares
	// 排序
	requestMiddlewareList := utils.SortMapByValue(requestMiddleware)

	// 根据排序结果
	for _, v := range requestMiddlewareList {
		middlewareName := strings.Split(v.Key, ".")[1]
		middlewares = append(middlewares, trackerMiddleware[middlewareName])
	}
	manager.MiddlewaresMap.Store(mapName, middlewares)
	return middlewares, nil
}

// NewMiddlewareManager 返回中间件管理器
func NewMiddlewareManager() *Manager {
	return &Manager{
		MiddlewaresMap: sync.Map{},
	}
}
