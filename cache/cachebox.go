package cache

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "httpproxy/lib"

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
    log.Println("yes to redis")
    return &CacheBox{
        pool: pool,
    }

}

func (c *CacheBox) Get(uri string) lib.Cache {
    log.Println("get cahche of ", uri)
    if cache := c.get(uri); cache != nil {
        //log.Println(*cache)
        return cache
    }
    return nil
}

func (c *CacheBox) get(uri string) *Cache {
    conn := c.pool.Get()
    defer conn.Close()

    b, err := redis.Bytes(conn.Do("GET", uri))
    if err != nil || len(b) == 0 {
        log.Println(err)
        return nil
    }
    log.Println(string(b))
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

func (c *CacheBox) CheckAndStore(uri string, resp *http.Response) {
    if !IsCache(resp) {
        return
    }

    cache := New(resp)

    if cache == nil {
        return
    }

    log.Println("store cache ", uri)
    b, err := json.Marshal(cache)
    if err != nil {
        log.Println(err)
        return
    }

    conn := c.pool.Get()
    defer conn.Close()

    conn.Send("MULTI")
    conn.Send("SET", uri, b)
    conn.Send("EXPIRE", uri, cache.maxAge)
    _, err = conn.Do("EXEC")
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("successfully store cache ", uri)

}

func (c *CacheBox) Clear(d time.Duration) {

}
