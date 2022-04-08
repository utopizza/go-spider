/* spider.go - spider爬虫对象封装 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
spider爬虫对象封装
*/

// Package spider 爬虫对象
package spider

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/utopizza/go-spider/cache"
	"github.com/utopizza/go-spider/configs"
	"github.com/utopizza/go-spider/parser"
	"github.com/utopizza/go-spider/storage"
)

// Service 爬虫服务接口
type Service interface {
	Start(seeds []string) error
}

// service 爬虫服务实现
type service struct {
	config         configs.Config
	cacheService   cache.Service
	parserService  parser.Service
	storageService storage.Service
	tokens         chan struct{}   // 控制并发数
	wg             *sync.WaitGroup // 主go程退出
	httpClient     *http.Client
}

// NewSpiderService 根据配置文件路径参数，创建爬虫对象，完成初始化
func NewSpiderService(config *configs.Config, cacheService cache.Service,
	parserService parser.Service, storageService storage.Service) (Service, error) {
	if cacheService == nil || parserService == nil || storageService == nil || config == nil {
		return nil, errors.New("nil input param")
	}
	s := new(service)
	s.config = *config
	s.cacheService = cacheService
	s.parserService = parserService
	s.storageService = storageService
	s.tokens = make(chan struct{}, config.Spider.ThreadCount)
	s.wg = &sync.WaitGroup{}
	s.httpClient = &http.Client{
		Timeout: time.Second * time.Duration(config.Spider.CrawlTimeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return s, nil
}

// Run 运行爬虫
func (s *service) Start(seeds []string) error {
	log.Info("start crawling webpages")
	for _, url := range seeds {
		s.wg.Add(1)
		go s.crawl(url, 0)
	}
	s.wg.Wait()
	log.Info("finish crawling webpages")
	return nil
}

// crawl 对目标url执行爬取任务
func (s *service) crawl(url string, depth int) {
	log.Infof("crawling webpage:%s, current depth:%d", url, depth)
	defer s.wg.Done()

	// 爬取深度控制
	if depth > s.config.Spider.MaxDepth {
		return
	}

	// 爬取时间间隔控制
	time.Sleep(time.Duration(s.config.Spider.CrawlInterval) * time.Second)

	// 爬取页面，提取新url
	s.tokens <- struct{}{}
	newURLs, err := s.extractURLs(url)
	if err != nil {
		log.Error(err)
	}
	<-s.tokens

	// 递归爬取新url
	for _, url := range newURLs {
		// 新url去重存储
		if !s.cacheService.PutIfAbsent(url) {
			log.Debugf("discard duplicated url:%s", url)
			continue
		}
		// 开启新go程进行爬取
		s.wg.Add(1)
		go s.crawl(url, depth+1)
	}
}

// extractURLs 从目标url对应的网页中提取新url
func (s *service) extractURLs(url string) ([]string, error) {
	// 请求页面
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code:%d of url:%s", resp.StatusCode, url)
	}

	// 从body中读出bytes
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 保存页面
	if err := s.storageService.Save(url, bodyBytes); err != nil {
		return nil, err
	}

	// 解析页面
	newURLs, err := s.parserService.Parse(resp, bodyBytes)
	if err != nil {
		return nil, err
	}

	return newURLs, nil
}
