/* main.go - 主程序 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
主程序入口，解析启动命令flags后，创建并运行spider对象爬取网页
*/

// Package main 主程序
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/utopizza/go-spider/cache"
	"github.com/utopizza/go-spider/configs"
	"github.com/utopizza/go-spider/parser"
	"github.com/utopizza/go-spider/spider"
	"github.com/utopizza/go-spider/storage"
)

const (
	version = "v0.1.0"
)

func main() {
	// 处理flag
	hPtr := flag.Bool("h", false, "show help")
	vPtr := flag.Bool("v", false, "version")
	cPtr := flag.String("c", "", "config directory")
	lPtr := flag.String("l", "", "log directory")

	flag.Parse()

	// flag参数检查
	if *hPtr {
		flag.PrintDefaults()
		return
	}
	if *vPtr {
		fmt.Printf("go-spider version: %s\n", version)
		return
	}
	if *lPtr == "" {
		fmt.Println("please use -l to specify log directory")
		flag.PrintDefaults()
		return
	}
	if *cPtr == "" {
		fmt.Println("please use -c to specify config directory")
		flag.PrintDefaults()
		return
	}
	configFile := path.Join(*cPtr, "spider.conf")

	// 初始化logger
	dir, err := filepath.Abs(*lPtr)
	if err != nil {
		fmt.Printf("illegal directory:%s", *lPtr)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.FileMode((0777)))
		if err != nil {
			fmt.Printf("fail to create directory:%s", dir)
		}
	}
	path := filepath.Join(dir, "log.txt")
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("fail to create file:%s", path)
	}
	defer file.Close()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)
	log.SetOutput(file)

	// 读取config
	config, err := configs.LoadConfig(configFile)
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 创建cache服务
	cacheService, err := cache.NewCacheService()
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 创建parser服务
	parserService, err := parser.NewParserService()
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 创建storage服务
	storageService, err := storage.NewStorageService(config.Spider.OutputDirectory, config.Spider.TargetURLPattern)
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 创建spider服务
	spiderService, err := spider.NewSpiderService(config, cacheService, parserService, storageService)
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 读取种子
	seeds, err := configs.LoadSeedFile(config.Spider.URLListFile)
	if err != nil {
		log.Error(err)
		exit(1)
	}

	// 开始爬取
	err = spiderService.Start(seeds)
	if err != nil {
		log.Error(err)
		exit(1)
	}
}

// exit main函数退出统一方法
func exit(code int) {
	os.Exit(code)
}
