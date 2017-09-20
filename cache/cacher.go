package cache

import (
	"net/http"
	"time"
)

// Box stores http get response cache
type Box interface {
	Get(uri string) Item
	Delete(uri string)
	CheckAndStore(uri string, resp *http.Response)
	Clear(d time.Duration)
}

// Item represent a http get response cache
type Item interface {
	Verify() bool
	WriteTo(rw http.ResponseWriter) (int, error)
}
