/* webpage_parse.go - 存储爬取到的网页内容 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
存储爬取到的网页内容
*/

// Package storage 网页存储服务
package storage

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
)

// Service 网页存储服务接口
type Service interface {
	Save(urlStr string, bytes []byte) error
}

// service 网页存储服务器
type service struct {
	dir string
	reg *regexp.Regexp
}

// NewStorageService 创建网页存储服务器
func NewStorageService(outputDir, targetURLPattern string) (Service, error) {
	s := new(service)
	s.dir = outputDir
	r, err := regexp.Compile(targetURLPattern)
	if err != nil {
		return nil, err
	}
	s.reg = r
	return s, nil
}

// Save 保存网页
func (s *service) Save(urlStr string, bytes []byte) error {
	// 检查目录
	if err := s.checkDir(); err != nil {
		return err
	}

	// 判断是否符合目标pattern
	if !s.reg.MatchString(urlStr) {
		return nil
	}

	// 以url为文件名，并escape特殊字符
	fileName := url.PathEscape(urlStr)

	// 创建文件
	f, err := os.Create(path.Join(s.dir, fileName))
	if err != nil {
		return fmt.Errorf("fail to create local file, url:%s, err:%v", urlStr, err)
	}
	defer f.Close()

	// 写字节数组到文件
	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("fail to save webpage, url:%s, err:%v", urlStr, err)
	}

	return nil
}

// checkDir 确认目录存在，否则创建目录。若创建失败则返回错误
func (s *service) checkDir() error {
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		if err = os.Mkdir(s.dir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to create output dir:%s, err:%+v", s.dir, err)
		}
	}
	return nil
}
