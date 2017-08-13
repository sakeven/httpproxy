package redis

import "github.com/sakeven/httpproxy/cache"

func Register(address string, password string) {
	cache.Register(NewCacheBox(address, password))
}
