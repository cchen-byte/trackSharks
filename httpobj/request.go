package httpobj

import (
	"fmt"
	"net/url"
)

// Request 请求对象
type Request struct {
	Url     string
	Method  string
	Params  map[string]interface{} // get 请求参数
	Data    map[string]interface{} // post 请求参数
	Payload map[string]string // Request Payload
	Json    interface{} // Request Json
	Headers map[string]string
	Timeout int    // 超时时间, 单位秒
	Ja3     string // ja3指纹
	Proxies string // scheme://username:password@ip:port

	Priority int	// 优先级

	DownloadMiddlewares map[string]int	// 下载中间件

	Export *Export // 需要导出的字段, 配置化/多步请求依赖前几步请求结果时 会用上

	IsRetry bool	// 是否重试
	RetryTimes int // 重试次数

	Callback func(response *Response) (*ParseResult, error) // 回调函数
}

func (request *Request) ToValues(args map[string]interface{}) string {
	if args != nil && len(args) > 0 {
		params := url.Values{}
		for k, v := range args {
			params.Set(k, fmt.Sprintf("%v", v))
		}
		return params.Encode()
	}
	return ""
}


// Export 导出的字段
type Export map[string]interface{}
