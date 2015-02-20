//Package cache handlers http web cache.
package cache

import (
	"crypto/sha1"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

// Copyheaders copys headers from response headers.
func (c *Cache) CopyHeaders(src http.Header) {
	c.Lock()
	defer c.Unlock()

	c.Header = make(http.Header)
	for key, values := range src {
		for _, value := range values {
			c.Header.Add(key, value)
		}
	}
}

//SetCache sets a new cache.
func (c *Cache) SetCache(StatusCode int, Body string) (err error) {
	c.Lock()
	defer c.Unlock()

	c.StatusCode = StatusCode
	c.Body = Body

	if c.Header == nil {
		return errors.New("try to access nil Header of a cache!")
	}

	c.ETag = c.Header.Get("ETag")
	c.Last_Modified = c.Header.Get("Last-Modified")

	Cache_Control := c.Header.Get("Cache-Control")
	if strings.Index(Cache_Control, "no-cache") != -1 {
		c.Mustverified = true
		return nil
	}

	if Expires := c.Header.Get("Expires"); Expires != "" {
		c.Vlidity, err = time.Parse(http.TimeFormat, Expires)
		if err != nil {
			return
		}
	}
	max_age := getAge(Cache_Control)
	if max_age != -1 {
		var Time time.Time
		date := c.Header.Get("Date")
		if date == "" {
			Time = time.Now().UTC()
		} else {
			Time, err = time.Parse(time.RFC1123, date)
			if err != nil {
				return
			}
		}
		c.Vlidity = Time.Add(time.Duration(max_age) * time.Second)
	}
	return nil
}

// Verify verifies whether cache is out of date.
func (c *Cache) Verify() bool {
	c.Lock()
	defer c.Unlock()

	newReq, err := http.NewRequest("GET", c.URI, nil)
	if err != nil {
		return false
	}

	if c.Last_Modified != "" {
		newReq.Header.Add("If-Modified-Since", c.Last_Modified)
	}
	if c.ETag != "" {
		newReq.Header.Add("If-None-Match", c.ETag)
	}
	Tr := &http.Transport{Proxy: http.ProxyFromEnvironment}
	resp, err := Tr.RoundTrip(newReq)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusNotModified {
		return false
	}
	return true
}

type Checksum [20]byte
type CacheSet map[Checksum]*Cache

func getCheckSum(URI string) Checksum {
	return sha1.Sum([]byte(URI))
}

// GetCache finds specific cache dertermined by URI,if not Found nil will be return.
func (c *CacheSet) GetCache(URI string) *Cache {
	return (*c)[getCheckSum(URI)]
}

// Delete deletes specific cache.
func (c *CacheSet) Delete(URI string) {
	delete(*c, getCheckSum(URI))
}
func (c *CacheSet) DeleteByCheckSum(key Checksum) {
	delete(*c, key)
}

// New returns a new cache.
func (c *CacheSet) New(URI string) *Cache {
	(*c)[getCheckSum(URI)].URI = URI
	return (*c)[getCheckSum(URI)]
}

//getAge from Cache Control get cache's lifetime.
func getAge(Cache_Control string) (age int64) {
	f := func(sage string) int64 {
		var tmpAge int64
		idx := strings.Index(Cache_Control, sage)
		if idx != -1 {
			for i := idx + len(sage) + 1; i < len(Cache_Control); i++ {
				if Cache_Control[i] >= '0' && Cache_Control[i] <= '9' {
					tmpAge = tmpAge*10 + int64(Cache_Control[i])
				} else {
					break
				}
			}
			return tmpAge
		}
		return -1
	}
	if s_maxage := f("s-maxage"); s_maxage != -1 {
		return s_maxage
	}
	return f("max-age")
}
