/* webpage_save_test.go - 单元测试 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
单元测试
*/

// Package storage 网页存储服务
package storage

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"
)

func TestSaveWebpage(t *testing.T) {
	// target url pattern
	targetURLPattern := ".*.(htm|html)$"

	// target url
	targetURL := "http://www.baidu.com/index.html"

	// output dir
	outputDir := "./"

	// 获取response
	resp, err := http.Get(targetURL)
	if err != nil {
		t.Error(err)
	}

	// 创建网页存储服务器
	storager, err := NewStorageService(outputDir, targetURLPattern)
	if err != nil {
		t.Error(err)
	}

	// 测试
	// 从body中读出bytes
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if err := storager.Save(targetURL, bodyBytes); err != nil {
		t.Error(err)
	}

	// 清除临时生成的文件
	fileName := url.PathEscape(resp.Request.URL.String())
	savePath := path.Join(outputDir, fileName)
	os.Remove(savePath)
}
