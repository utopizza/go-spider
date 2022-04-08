/* webpage_parse.go - 解析网页内容 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
解析网页内容
*/

// Package parser 网页解析器
package parser

import (
	"bytes"
	"net/http"

	"golang.org/x/net/html"
)

// Service 网页解析服务接口
type Service interface {
	Parse(resp *http.Response, bodyBytes []byte) ([]string, error)
}

// service 网页解析器
type service struct {
}

// NewParserService 创建网页解析器
func NewParserService() (Service, error) {
	s := new(service)
	return s, nil
}

// Parse 解析网页，提取url
func (s *service) Parse(resp *http.Response, bodyBytes []byte) ([]string, error) {

	// 解析body
	entryNode, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	// 从页面提取url
	var links []string
	onEachNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {

				// 过滤出链接
				if a.Key != "href" {
					continue
				}

				// url统一处理为绝对路径
				newURL, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				newURLStr := newURL.String()
				links = append(links, newURLStr)
			}
		}
	}
	recursive(entryNode, onEachNode)

	return links, nil
}

// recursive 递归遍历网页中的节点进行解析
func recursive(node *html.Node, visitNodeFunc func(n *html.Node)) {
	if node == nil || visitNodeFunc == nil {
		return
	}
	visitNodeFunc(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		recursive(child, visitNodeFunc)
	}
}
