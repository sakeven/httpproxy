package cache

import (
    "encoding/json"
    // "log"
    "net/http"
    "time"

    "github.com/garyburd/redigo/redis"
)

type CacheBox struct {
    pool *redis.Pool
}

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

func (c *CacheBox) Get(uri string) *Cache {
    return c.get(uri)
}

func (c *CacheBox) get(uri string) *Cache {
    conn := c.pool.Get()
    defer conn.Close()

    values, err := redis.Values(conn.Do("GET", uri))

    var b []byte
    _, err = redis.Scan(values, &b)
    if err != nil {
        return nil
    }

    cache := new(Cache)
    json.Unmarshal(b, &cache)
    return cache
}

func (c *CacheBox) Delete(uri string) {
    c.delete(uri)
}

func (c *CacheBox) delete(uri string) {
    conn := c.pool.Get()
    defer conn.Close()

    _, err := conn.Do("DEL", uri)

    if err != nil {
        return
    }

    return
}

func (c *CacheBox) CheckAndStore(resp *http.Response) {
    if !IsCache(resp) {
        return
    }

    cache := New(resp)

    b, err := json.Marshal(cache)
    if err != nil {
        return
    }

    conn := c.pool.Get()
    defer conn.Close()

    conn.Send("MULTI")
    conn.Send("SET", resp.Request.RequestURI, b)
    conn.Send("EXPIRE", cache.StatusCode)
    _, err = conn.Do("EXEC")
    if err != nil {
        return
    }

}

func (c *CacheBox) CheckAndDelete(d time.Duration) {

}
