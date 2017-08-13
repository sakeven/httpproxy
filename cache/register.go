package cache

import (
	"net/http"
	"time"
)

var globalBox CacheBox

func Register(box CacheBox) {
	globalBox = box
}

func Get(uri string) Cache {
	return globalBox.Get(uri)
}

func Delete(uri string) {
	globalBox.Delete(uri)
}

func CheckAndStore(uri string, resp *http.Response) {
	globalBox.CheckAndStore(uri, resp)
}

func Clear(d time.Duration) {
	globalBox.Clear(d)
}
