/* webpage_parse_test.go - 单元测试 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
单元测试
*/

// Package parser 网页解析器
package parser

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	// 创建parser对象
	parser, err := NewParserService()
	if err != nil {
		t.Error(err)
	}

	// 获取response
	resp, err := http.Get("https://www.sina.com.cn/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	// 测试
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	newURLs, err := parser.Parse(resp, bodyBytes)
	if err != nil {
		t.Error(err)
	}
	if len(newURLs) == 0 {
		t.Error("fail to parse webpage")
	}
}
