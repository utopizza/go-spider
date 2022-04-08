/* url_table_test.go - 单元测试 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
单元测试
*/

// Package cache 缓存层
package cache

import "testing"

func TestNewCacheService(t *testing.T) {
	s, err := NewCacheService()
	if err != nil {
		t.Error(err)
	}
	if s == nil {
		t.Error("fail to NewCacheService")
	}
}

func TestPutIfAbsent(t *testing.T) {
	s, err := NewCacheService()
	if err != nil {
		t.Error(err)
	}
	targetURL := "https://www.baidu.com"
	if !s.PutIfAbsent(targetURL) {
		t.Error("fail to put new url")
	}
	if s.PutIfAbsent(targetURL) {
		t.Error("put duplicate url")
	}
}
