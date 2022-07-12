package httpobj

import (
	"bytes"
	"fmt"
	"github.com/cchen-byte/trackeSharkes/httpobj/dom"
)

// Response 响应对象
type Response struct {
	Request    *Request
	StatusCode int
	Url        string
	Headers    map[string][]string
	Content    []byte
	Text       string
}

// GetDom 获取dom
func (response *Response) GetDom() (*dom.Dom, error) {
	if len(response.Content) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewDom(bytes.NewReader(response.Content))
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetXmlDom 获取xmlDom
func (response *Response) GetXmlDom() (*dom.XmlDom, error) {
	if len(response.Content) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewXmlDom(bytes.NewReader(response.Content))
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetJsonDom 获取xmlDom
func (response *Response) GetJsonDom() (*dom.JsonDom, error) {
	if len(response.Content) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	d, err := dom.NewJsonDom(bytes.NewReader(response.Content))
	if err != nil {
		return nil, err
	}
	return d, nil
}


// ParseResult 回调函数解析的结果
type ParseResult struct {
	Requests []*Request // 解析出来的请求
	Items    []*Item    // 解析出来的内容
}

// AppendItem 追加item
func (pr *ParseResult) AppendItem(item *Item) {
	pr.Items = append(pr.Items, item)
}

// AppendRequest 追加Request
func (pr *ParseResult) AppendRequest(req *Request) {
	pr.Requests = append(pr.Requests, req)
}

func NewParseResult() *ParseResult {
	return &ParseResult{}
}