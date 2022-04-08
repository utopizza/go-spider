/* load_config.go - 读取配置文件 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
根据文件路径读取配置
*/

// Package configs 配置层
package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/gcfg.v1"
)

// Config 配置汇总
type Config struct {
	Spider SpiderConfig
}

// check 配置汇总检查
func (c *Config) check() error {
	if err := c.Spider.check(); err != nil {
		return err
	}
	return nil
}

// SpiderConfig spider配置
type SpiderConfig struct {
	URLListFile      string `gcfg:"urlListFile"`
	OutputDirectory  string `gcfg:"outputDirectory"`
	MaxDepth         int    `gcfg:"maxDepth"`
	CrawlInterval    int    `gcfg:"crawlInterval"`
	CrawlTimeout     int    `gcfg:"crawlTimeout"`
	TargetURLPattern string `gcfg:"targetUrl"`
	ThreadCount      int    `gcfg:"threadCount"`
}

// check spider配置内容检查，例如非空检查、数值范围检查、正则表达式合法性检查等
func (c *SpiderConfig) check() error {
	if c.URLListFile == "" {
		return errors.New("empty urlListFile")
	}
	if c.OutputDirectory == "" {
		return errors.New("empty outputDirectory")
	}
	if c.MaxDepth < 0 {
		return errors.New("maxDepth must be positve integer")
	}
	if c.CrawlInterval < 0 {
		return errors.New("crawlInterval must be postive integer")
	}
	if c.CrawlTimeout < 0 {
		return errors.New("crawlTimeout must be postive integer")
	}
	if c.ThreadCount < 0 {
		return errors.New("threadCount must be postive integer")
	}
	_, err := regexp.Compile(c.TargetURLPattern)
	if err != nil {
		return fmt.Errorf("illegal targetURLPattern, err:%v", err.Error())
	}
	return nil
}

// LoadConfig 加载配置
func LoadConfig(path string) (*Config, error) {
	c := new(Config)
	err := gcfg.ReadFileInto(c, path)
	if err != nil {
		return nil, err
	}

	// 配置内容检查
	if err := c.check(); err != nil {
		return nil, err
	}

	return c, nil
}

// LoadSeedFile 根据文件路径读取种子URL
func LoadSeedFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var seeds []string
	err = json.Unmarshal(bytes, &seeds)
	if err != nil {
		return nil, err
	}

	return seeds, nil
}
