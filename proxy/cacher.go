package proxy

import (
    "net/http"
    "time"
)

type CacheBox interface {
    Get(uri string) Cache
    Delete(uri string)
    CheckAndStore(resp *http.Response)
    CheckAndDelete(d time.Duration)
}

type Cache interface {
    Verify() bool
    WriteTo(rw http.ResponseWriter) (int, error)
}
