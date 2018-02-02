## httpproxy
[![Go Report Card](https://goreportcard.com/badge/github.com/sakeven/httpproxy)](https://goreportcard.com/report/github.com/sakeven/httpproxy)

* 基于 Go 开发，支持 HTTP/1.1 以上版本的 HTTP 代理。

## 细节功能：
* 支持内容缓存和重校验
* 支持 GET、POST、CONNECT 等方法
* 支持账号登入与验证
* 支持配置文件
* 提供 Web 版管理和调试界面
* 支持反向代理

## 正在进行中
* 资源限定(各种超时，最大文件大小，最大缓存大小，最大头大小等，最大并发量，最大请求速度，最大传输速度等)

## 配置
  
配置文件在 $HOME/.httproxy/config.json，采用 JSON 格式，包含

* port：代理服务器工作端口
* webport：代理服务器 web 管理端口
* reverse：设置反向代理，值为 true 或者 false
* proxy_pass：反向代理目标服务器地址，如 "127.0.0.1:80"
* auth：开启代理认证，值为 true 或者 false
* cache：开启缓存，值为 true 或者 false
* cache_timeout：缓存更新时间，单位分钟
* log：值为 1 时输出调试信息，为 0 时输出普通监控信息
* gfwlist：网站屏蔽列表，如 ["baidu.com","google.com"]
* admin：web 管理用户
* user：代理服务器普通用户

一个简单配置演示如下

```json
{
    "port": ":8080",
    "webport": ":6060",
    "reverse": true,
    "proxy_pass": "127.0.0.1:80",
    "auth": true,
    "cache": false,
    "cache_timeout": 30,
    "log": 1,
    "gfwlist": [
        "baidu.com",
        "google.com"
    ],
    "admin": {
        "Admin": "prxy"
    },
    "user": {
        "proxy": "proxy"
    }
}
/*this is a configure for proxy server. log: 1 for Information, 0 for DebugInfor*/
```

## 构建
* 安装

        $ go get github.com/sakeven/httpproxy

* 运行

        $ ./${GOPATH}/bin/httpproxy

