/* load_config_test.go - 单元测试 */
/*
modification history
--------------------
2021/01/14, by wangyusheng01, create
*/
/*
DESCRIPTION
单元测试
*/

// Package configs 配置层
package configs

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("../conf/spider.conf")
	if err != nil {
		t.Error(err)
	}
	if c.Spider.MaxDepth != 1 {
		t.Error("fail to load config file")
	}
	_, err = LoadConfig("")
	if err == nil {
		t.Error("load empty file")
	}
}

func TestLoadSeedFile(t *testing.T) {
	seeds, err := LoadSeedFile("../data/url.data")
	if err != nil {
		t.Error(err)
	}
	if len(seeds) == 0 {
		t.Error("fail to load seed file")
	}
	_, err = LoadSeedFile("")
	if err == nil {
		t.Error("load empty file")
	}
}
