package redis

import (
	"net/http"
	"strings"
)

// IsCache checks whether response can be stored as cache
func IsCache(resp *http.Response) bool {
	cacheControl := resp.Header.Get("Cache-Control")
	contentType := resp.Header.Get("Content-Type")
	if strings.Index(cacheControl, "private") != -1 ||
		strings.Index(cacheControl, "no-store") != -1 ||
		strings.Index(contentType, "application") != -1 ||
		strings.Index(contentType, "video") != -1 ||
		strings.Index(contentType, "audio") != -1 ||
		(strings.Index(cacheControl, "max-age") == -1 &&
			strings.Index(cacheControl, "s-maxage") == -1 &&
			resp.Header.Get("Etag") == "" &&
			resp.Header.Get("Last-Modified") == "" &&
			(resp.Header.Get("Expires") == "" || resp.Header.Get("Expires") == "0")) {
		return false
	}
	return true
}

// //CheckCaches evey certain minutes check whether cache is out of date, if yes release it.
// func CheckCaches() {
//     for {
//         time.Sleep(time.Duration(cnfg.CacheTimeout) * time.Minute)
//         for key, Cache := range Caches {
//             if Cache != nil && Cache.Verify() == false {
//                 Caches.DeleteByCheckSum(key)
//             }
//         }
//     }
// }
