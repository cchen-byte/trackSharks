package engine

import "github.com/cchen-byte/trackeSharkes/httpobj"

type Engine interface {
	SubmitRequests(*httpobj.Request) // 提交一个 Request
	SubmitParseResult(result *httpobj.ParseResult)	// 提交一个 ParseResult
	Run()
}

//type GetEngineFactory interface {
//	GetEngine() Engine
//}