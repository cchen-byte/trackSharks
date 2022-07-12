package dom

import (
	"github.com/andybalholm/cascadia"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io"
)

type Dom struct {
	node      *html.Node
	Type      html.NodeType
	Data      string
	Namespace string
	Attr      []html.Attribute
}

// RawNode 返回原始html.node
func (d *Dom) RawNode() *html.Node {
	return d.node
}

// Xpath xpath查找所有, 表达式错误会panic
func (d *Dom) Xpath(expr string) []*Dom {
	nodes := htmlquery.Find(d.node, expr)
	return makeDom(nodes)
}

// XpathOne xpath查找一个, 表达式错误会panic
func (d *Dom) XpathOne(expr string) *Dom {
	node := htmlquery.FindOne(d.node, expr)
	return makeOneDom(node)
}

// Css css查找所有, 表达式错误会panic
func (d *Dom) Css(expr string) []*Dom {
	sel := getCssQuery(expr)
	nodes := cascadia.QueryAll(d.node, sel)
	return makeDom(nodes)
}

// CssOne css查找一个, 表达式错误会panic
func (d *Dom) CssOne(expr string) *Dom {
	sel := getCssQuery(expr)
	node := cascadia.Query(d.node, sel)
	return makeOneDom(node)
}

// GetAttr 获取属性
func (d *Dom) GetAttr(name string) string {
	return htmlquery.SelectAttr(d.node, name)
}

// InnerText 获取node内所有的文本值
func (d *Dom) InnerText() string {
	return htmlquery.InnerText(d.node)
}

/* 将node转为html
   self 表示是否输入自己
*/
func (d *Dom) OutputHTML(self bool) string {
	return htmlquery.OutputHTML(d.node, self)
}

// Parent 返回上级节点
func (d *Dom) Parent() *Dom {
	return makeOneDom(d.node.Parent)
}

// FirstChild 返回第一个子节点
func (d *Dom) FirstChild() *Dom {
	return makeOneDom(d.node.FirstChild)
}

// LastChild 返回最后一个子节点
func (d *Dom) LastChild() *Dom {
	return makeOneDom(d.node.LastChild)
}

// PrevSibling 返回上一个同级节点
func (d *Dom) PrevSibling() *Dom {
	return makeOneDom(d.node.PrevSibling)
}

// NextSibling 返回下一个同级节点
func (d *Dom) NextSibling() *Dom {
	return makeOneDom(d.node.NextSibling)
}

// Children 获取所有子
func (d *Dom) Children() []*Dom {
	var a []*html.Node
	for nn := d.node.FirstChild; nn != nil; nn = nn.NextSibling {
		a = append(a, nn)
	}
	return makeDom(a)
}

func NewDom(r io.Reader) (*Dom, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return makeOneDom(node), nil
}

func makeDom(nodes []*html.Node) []*Dom {
	dom := make([]*Dom, len(nodes))
	for i, node := range nodes {
		dom[i] = makeOneDom(node)
	}
	return dom
}

func makeOneDom(node *html.Node) *Dom {
	if node == nil {
		return nil
	}
	return &Dom{
		node:      node,
		Type:      node.Type,
		Data:      node.Data,
		Namespace: node.Namespace,
		Attr:      node.Attr,
	}
}
