package middleware

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/setting"
	mapset "github.com/deckarep/golang-set"
)

// RetryMiddleware 内置重试中间件
type RetryMiddleware struct {

}

func (rm *RetryMiddleware) ProcessRequest(request *httpobj.Request) error {
	return nil
}


func (rm *RetryMiddleware) ProcessResponse(response *httpobj.Response) error {
	statusCodeList := []interface{}{200}
	statusCode := mapset.NewSetFromSlice(statusCodeList)
	// 当前 response.StatusCode 不在 statusCodeList 内 && response 的重试次数小于设置的次数
	if !statusCode.Contains(response.StatusCode) && response.Request.RetryTimes < setting.Settings.RetryTime {
		response.Request.RetryTimes += 1
		response.Request.IsRetry = true
	}else{
		response.Request.IsRetry = false
	}
	return nil
}