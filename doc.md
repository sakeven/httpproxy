#包名：cache
导入方式：`import "httpproxy/cache"`  
包文件：httpproxy/cache.go  
功能：该包处理了代理服务器的缓存，可以将缓存储存在内存中，支持对缓存的验证。
##type Cache
	type Cache struct {
	   Header        http.Header
	   Body          string
	   StatusCode    int
	   URI           string
	   Last_Modified string //eg:"Fri, 27 Jun 2014 07:19:49 GMT"
	   ETag          string
	   Mustverified  bool
	   //Vlidity is a time when to verfiy the cache again.
	   Vlidity time.Time
	   sync.Mutex
	}

`func (c *Cache) CopyHeaders(src http.Header)`  
`CopyHeaders`方法复制http响应头到`Cache`中

`func (c *Cache) SetCache(StatusCode int, Body string) (err error)`  
`SetCache`方法通过`StatusCode`状态码，`Body`响应主体设置一个新的`Cache`，设置失败时，将返回一个`error`。

`func (c *Cache) Verify() bool`  
`Verify`方法验证缓存`c`是否过期，过期返回`false`，未过期返回`true`。

##type CacheSet
	type CacheSet map[string]*Cache
`CacheSet` 是`Cache`的集合

`func (c *CacheSet) Delete(URI string)`  
`Delete`方法从缓存集合`c`中删除由`URI`指定的`Cache`


`func (c *CacheSet) GetCache(URI string) *Cache`  
`GetCache`方法查找`URI`指定的缓存`cache`，如果没有找到将返回一个`nil`。

#包名 config
导入方式：`import "httpproxy/config"`  
包文件：httpproxy/config.go  
功能：包config提供对json格式配置文件的读取和保存。  

## type Config

	type Config struct {
	    // 代理服务器工作端口,eg:":8080"
	    Port string `json:"port"`
	
	    // web管理端口
	    WebPort string `json:"webport"`
	
	    // 反向代理标志
	    Reverse bool `json:"reverse"`
	
	    // 反向代理目标地址,eg:"127.0.0.1:8090"
	    Proxy_pass string `json:"proxy_pass"`
	
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
`Config` 保存代理服务器的配置

`func (c *Config) GetConfig(filename string) error`  
`GetConfig`方法从指定json文件读取config配置

`func (c *Config) WriteToFile(filename string) error`  
`WriteToFile`方法将config配置写入特定json文件

# 包 proxy
导入方式： `import "httpproxy/proxy"`  
包文件：httpproxy/proxy/auth.go httpproxy/proxy/cache.go httpproxy/proxy/headers.go httpproxy/proxy/init.go httpproxy/proxy/proxy.go httpproxy/proxy/web.go  
功能：Package proxy 实现了一个http代理服务器，支持GET,POST,CONNECT等方法，支持代理认证和web管理，支持缓存和反向代理。

##全局变量：

`var Caches cache.CacheSet`缓存容器

`var HTTP_200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")`

`var HTTP_401 = []byte("HTTP/1.1 401 Authorization Required\r\nWWW-Authenticate: Basic realm=\"Secure Web\"\r\n\r\n")`

`var HTTP_407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Secure Proxys\"\r\n\r\n")`

##函数：
`func Check(user, passwd string) bool`  
`check`函数验证用户名和密码

`func CheckAdmin(user, passwd string) bool`
`CheckAdmin`函数验证admin管理，用于web管理界面登入

`func CheckCaches()`
`CheckCaches`函数每隔特定时间释放过期缓存

`func ClearHeaders(headers http.Header)`  
`ClearHeaders`函数清空HTTP头

`func CopyHeaders(rw, resp http.Header)`  
`CopyHeaders`函数将HTTP头从resp拷贝到rw。

`func ExistCache(uri string) bool`  
`ExistCache`函数检查指定uri的缓存是否已存在

`func IsCache(resp *http.Response, URI string) bool`  
`IsCache`函数通过对http头的检测确定响应是否应该缓存

`func NeedAuth(rw http.ResponseWriter, challenge []byte) error`  
`NeedAuth`函数在代理服务器需要验证时，发送验证响应

`func NewProxyServer() *http.Server`  
`NewProxyServer`函数返回一个新的代理服务器

`func RmProxyHeaders(req *http.Request)`  
`RmProxyHeaders`函数删除http逐跳头

##type ProxyServer

	type ProxyServer struct {
	    // User records user's name
	    Tr   *http.Transport
	    User string
	}

`func (proxy *ProxyServer) Auth(rw http.ResponseWriter, req *http.Request) (string, error)`  
`Auth`方法验证用户登入，如果不能登入将返回一个error

`func (proxy *ProxyServer) CacheHandler(rw http.ResponseWriter, req *http.Request)`  
`CacheHandler`方法只处理get请求的缓存

`func (proxy *ProxyServer) HttpHandler(rw http.ResponseWriter, req *http.Request)`  
`HttpHandler`方法处理普通的http请求

`func (proxy *ProxyServer) HttpsHandler(rw http.ResponseWriter, req *http.Request)`  
`HttpsHandler`方法处理https连接，主要用于CONNECT方法

`func (proxy *ProxyServer) ReverseHandler(req *http.Request)`  
`ReverseHandle`方法处理反向代理请求

`func (proxy *ProxyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request)`  
`ServeHTTP` will be automatically called by system. `ProxyServer` implements the Handler interface which need ServeHTTP.

##type WebServer

	type WebServer struct {
	   //端口
	   Port string
	}

`func NewWebServer() *WebServer`  
`NewWebServer`函数返回一个新的web管理服务器

`func (ws *WebServer) HomeHandler(rw http.ResponseWriter, req *http.Request)`  
`HomeHandler`方法处理home页面

`func (ws *WebServer) ServeHTTP(rw http.ResponseWriter, req *http.Request)`  
`ServeHTTP`方法处理web管理页面的所有请求，并将请求转发到特定函数

`func (ws *WebServer) SettingHandler(rw http.ResponseWriter, req *http.Request)`  
`SettingHandler`方法用于代理服务器的配置文件设置

`func (ws *WebServer) UserHandler(rw http.ResponseWriter, req *http.Request)`  
`UserHandler`方法处理用户列表，有列出、修改、删除、增加用户等功能

`func (ws *WebServer) WebAuth(rw http.ResponseWriter, req *http.Request) error`  
`WebAuth`方法处理web管理页面的管理元验证与登入
