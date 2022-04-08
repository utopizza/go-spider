/* url_table.go - 存储已经爬取过的url */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
存储已经爬取过的url
*/

// Package cache 缓存层
package cache

import (
	"sync"
)

// Service 缓存服务接口
type Service interface {
	PutIfAbsent(url string) bool
}

// service 缓存服务实现
type service struct {
	cache map[string]bool
	mutex *sync.Mutex
}

// NewCacheService 创建缓存服务实例
func NewCacheService() (Service, error) {
	s := new(service)
	s.cache = make(map[string]bool)
	s.mutex = &sync.Mutex{}
	return s, nil
}

// PutIfAbsent 若缓存中不存在该url则进行缓存并返回true，否则返回false表示已在缓存中
func (s *service) PutIfAbsent(url string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cache[url] {
		return false
	}
	s.cache[url] = true
	return true
}
