package redis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	cache "github.com/sakeven/httpproxy/cache"
)

// MD5URI returns md5 hash of uri.
func MD5URI(uri string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(uri)))
}

// CacheBox is an redis instance to store cache.
type CacheBox struct {
	pool *redis.Pool
}

// NewCacheBox creates a redis cache implemented cache.Box.
func NewCacheBox(address string, password string) *CacheBox {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 1 * time.Hour,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err = c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}

	c := pool.Get()
	defer c.Close()

	_, err := c.Do("PING")
	if err != nil {
		panic("Fail to connect to redis server")
	}
	return &CacheBox{
		pool: pool,
	}

}

// Get gets an item of specific uri.
func (c *CacheBox) Get(uri string) cache.Item {
	log.Println("get cahche of ", uri)
	if cache := c.get(MD5URI(uri)); cache != nil {
		//log.Println(*cache)
		return cache
	}
	return nil
}

func (c *CacheBox) get(md5URI string) *Cache {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", md5URI))
	if err != nil || len(b) == 0 {
		return nil
	}
	cache := new(Cache)
	json.Unmarshal(b, &cache)
	return cache
}

// Delete deletes an item of specific uri.
func (c *CacheBox) Delete(uri string) {
	c.delete(MD5URI(uri))
}

func (c *CacheBox) delete(md5URI string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", md5URI)

	if err != nil {
		return
	}

	return
}

// CheckAndStore checks resp and then store it of specific uri.
func (c *CacheBox) CheckAndStore(uri string, resp *http.Response) {
	if !IsCache(resp) {
		return
	}

	cache := New(resp)

	if cache == nil {
		return
	}

	log.Println("store cache ", uri)

	md5URI := MD5URI(uri)
	b, err := json.Marshal(cache)
	if err != nil {
		log.Println(err)
		return
	}

	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SET", md5URI, b)
	conn.Send("EXPIRE", md5URI, cache.maxAge)
	_, err = conn.Do("EXEC")
	if err != nil {
		return
	}
	log.Println("successfully store cache ", uri)

}

// Clear refresh cache box in d.
func (c *CacheBox) Clear(d time.Duration) {}
