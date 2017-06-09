## httpproxy
* 基于golang开发，支持http/1.1以上版本的http代理。

## 细节功能：
* 支持内容缓存和重校验
* 支持GET\POST\CONNECT等方法
* 支持账号登入与验证
* 支持配置文件
* 提供web版管理和调试界面
* 支持反向代理

## 正在进行中
* 资源限定(各种超时，最大文件大小，最大缓存大小，最大头大小等，最大并发量，最大请求速度，最大传输速度等)

## 配置
  
配置文件在config目录，采用json格式，包含

* port：代理服务器工作端口
* webport：代理服务器web管理端口
* reverse：设置反向代理，值为true或者false
* proxy_pass：反向代理目标服务器地址，如"127.0.0.1:80"
* auth：开启代理认证，值为true或者false
* cache：开启缓存，值为true或者false
* cache_timeout：缓存更新时间，单位分钟
* log：值为1时输出Debug调试信息，为0时输出普通监控信息
* gfwlist：网站屏蔽列表，如["baidu.com","google.com"]
* admin：web管理用户
* user：代理服务器普通用户

一个简单配置演示如下

    {
        "port":":8080",
        "webport":":6060",
        "reverse":true,
        "proxy_pass":"127.0.0.1:80",
        "auth":true,
        "cache":false,
        "cache_timeout":30,
        "log":1,
        "gfwlist":["baidu.com","google.com"],
        "admin":{"Admin":"prxy"},
        "user":{"proxy":"proxy"}
    }


    /*this is a configure for proxy server. log: 1 for Information, 0 for DebugInfor*/

## Build
* 在$GOPATH/src目录

        $ git clone git://code.csdn.net/sakeven/httpproxy.git
* 安装go-logging

        $ go get github.com/op/go-logging
* 打开$GOPATH/src/httpproxy目录，并编译

        $ cd $GOPATH/src/httpporxy
        $ go build
* 运行

        $ ./httpproxy

## Bug 
Contact with the author jc5930@sina.cn
