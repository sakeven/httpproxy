// Package config provides Config struct for proxy.
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 保存代理服务器的配置
type Config struct {
	// 代理服务器工作端口,eg:":8080"
	Port string `json:"port"`

	// web管理端口
	WebPort string `json:"webport"`

	// 反向代理标志
	Reverse bool `json:"reverse"`

	// 反向代理目标地址,eg:"127.0.0.1:8090"
	ProxyPass string `json:"proxy_pass"`

	// 认证标志
	Auth bool `json:"auth"`

	// 缓存标志
	Cache bool `json:"cache"`

	// 缓存定期刷新时间，单位分钟
	CacheTimeout int64 `json:"cache_timeout"`

	// 日志信息，1输出Debug信息，0输出普通监控信息
	Log int `json:"log"`

	// 网站屏蔽列表
	GFWList []string `json:"gfwlist"`

	// 管理员账号
	Admin map[string]string `json:"admin"`
	// 普通用户账户
	User map[string]string `json:"user"`
}

var configFile = os.Getenv("HOME") + "/.httpproxy/config.json"

func isExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("file does not exist")
		return false
	}
	return true
}

func createDefaultConfig() {
	err := os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("create dir %s failed", filepath.Dir(configFile)))
	}

	file, err := os.Create(configFile)
	if err != nil {
		panic(fmt.Sprintf("open %s failed", configFile))
	}
	defer file.Close()

	file.Write([]byte(defaultConfig))
}

// GetConfig gets config from json file.
// GetConfig 从指定json文件读取config配置
func (c *Config) GetConfig() error {
	if isExist(configFile) == false {
		createDefaultConfig()
	}

	c.Admin = make(map[string]string)
	c.User = make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	return json.NewDecoder(br).Decode(c)
}

// WriteToFile writes config into json file.
// WriteToFile 将config配置写入特定json文件
func (c *Config) WriteToFile() error {
	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	return err
}

const defaultConfig = `
{
    "port": ":8080",
    "webport": ":6060",
    "reverse": false,
    "proxy_pass": "127.0.0.1:80",
    "auth": false,
    "cache": false,
    "cache_timeout": 60,
    "log": 0,
    "gfwlist": [
    ],
    "admin": {
        "Admin": "prxy"
    },
    "user": {
        "proxy": "proxy"
    }
}
`
