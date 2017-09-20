package redis

import "github.com/sakeven/httpproxy/cache"

// Register register redis cahce box as a global cache.Box.
func Register(address string, password string) {
	cache.Register(NewCacheBox(address, password))
}
