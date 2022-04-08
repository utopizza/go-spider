/* spider_test.go - 单元测试 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
单元测试
*/

// Package spider 爬虫对象
package spider

import (
	"log"
	"testing"

	"github.com/utopizza/go-spider/cache"
	"github.com/utopizza/go-spider/configs"
	"github.com/utopizza/go-spider/parser"
	"github.com/utopizza/go-spider/storage"
)

var spiderService Service

func TestMain(m *testing.M) {
	// 读取config
	config, err := configs.LoadConfig("../conf/spider.conf")
	if err != nil {
		log.Fatal(err)
	}
	config.Spider.MaxDepth = 0

	// 创建cache服务
	cacheService, err := cache.NewCacheService()
	if err != nil {
		log.Fatal(err)
	}

	// 创建parser服务
	parserService, err := parser.NewParserService()
	if err != nil {
		log.Fatal(err)
	}

	// 创建storage服务
	storageService, err := storage.NewStorageService(config.Spider.OutputDirectory, config.Spider.TargetURLPattern)
	if err != nil {
		log.Fatal(err)
	}

	// 创建spider
	spiderService, err = NewSpiderService(config, cacheService, parserService, storageService)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestNewSpiderService(t *testing.T) {
	_, err := NewSpiderService(nil, nil, nil, nil)
	if err == nil {
		t.Error("nil param to NewSpiderService")
	}
}

func TestStart(t *testing.T) {
	seeds := []string{"https://www.sina.com.cn/"}
	err := spiderService.Start(seeds)
	if err != nil {
		t.Error(err)
	}
}
