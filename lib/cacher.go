package lib

import (
    "net/http"
    "time"
)

type CacheBox interface {
    Get(uri string) Cache
    Delete(uri string)
    CheckAndStore(uri string, resp *http.Response)
    Clear(d time.Duration)
}

type Cache interface {
    Verify() bool
    WriteTo(rw http.ResponseWriter) (int, error)
}
