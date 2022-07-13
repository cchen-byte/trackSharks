package downloader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 在官方文档中有提到 net/http client 是协程安全的, 应该复用.
// 但在并发的过程中如果要更改重定向次数, 设置代理, 超时时间等的话, 就会有并发安全问题。
// 应当使用上下文的方式创建

// NetDownloader 原生http下载器
type NetDownloader struct {
	Session *Session
}

func (nd *NetDownloader) Fetch(request *httpobj.Request) (*httpobj.Response, error) {
	timeout := time.Second * 10
	// 定义超时时间, 单位: 秒
	if request.Timeout != 0 {
		timeout = time.Second * time.Duration(request.Timeout)
	}

	// 定义上下文
	ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)

	// 定义代理 set proxy to request context.
	if request.Proxies != "" {
		ctx = context.WithValue(ctx, "http", request.Proxies)
		ctx = context.WithValue(ctx, "https", request.Proxies)
	}

	if request.Method == "" {
		request.Method = "GET"
	}

	// Post Body
	var requestBody io.Reader

	// 请求参数处理
	switch strings.ToUpper(request.Method) {
	case "GET":

	case "POST":
		// 表单提交
		if request.Payload != nil {
			payload := make(url.Values)
			for key, value := range request.Payload{
				payload.Add(key, value)
			}
			requestBody = strings.NewReader(payload.Encode())
		}

		// json 提交
		if request.Json != nil {
			requestBodyByte, _ := json.Marshal(request.Json)
			requestBody = strings.NewReader(string(requestBodyByte))
		}

		// PostForm POST接口form表单
		if len(request.Data) > 0 {
			requestBody = strings.NewReader(request.ToValues(request.Data))
		}

	}
	// Get Params
	if len(request.Params) > 0 {
		request.Url = request.Url + "?" + request.ToValues(request.Params)
	}

	// 定义请求
	// Make new http.Request with context
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, requestBody)
	if err != nil {
		timeoutCancel()
		return nil, err
	}

	// 携带请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36")
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	//// 添加cookie
	//if rawCookie, ok := request.Headers["Cookie"]; ok {
	//	cookieList := ParseCookies(rawCookie)
	//	for _, cookie := range cookieList{
	//		req.AddCookie(cookie)
	//	}
	//}
	//delete(request.Headers, "Cookie")

	//
	resp, err := nd.Session.client.Do(req)
	if err != nil {
		timeoutCancel()
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	respByte, _ := ioutil.ReadAll(resp.Body)

	// response 封装
	responseData := &httpobj.Response{
		Request:    request,
		StatusCode: resp.StatusCode,
		Url:        request.Url,
		Headers:    resp.Header,
		Content:    respByte,
		Text:       string(respByte),
	}
	timeoutCancel()
	return responseData, nil
}

// 检测编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, e := r.Peek(1024)
	if e != nil {
		return unicode.UTF8
	}
	encode, _, _ := charset.DetermineEncoding(bytes, "")
	return encode
}

//
//type NetDownloaderFactory struct {}
//
//func (netDownloaderFactory *NetDownloaderFactory) GetDownloader() Downloader {
//	return &NetDownloader{
//		Session: NewSession(),
//	}
//}

func NewNetDownloader() *NetDownloader{
	return &NetDownloader{
		Session: NewSession(),
	}
}