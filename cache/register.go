package cache

import (
	"net/http"
	"time"
)

var globalBox Box

// Register registers a cache box in global.
func Register(box Box) {
	globalBox = box
}

// Get gets an item of specific uri.
func Get(uri string) Item {
	return globalBox.Get(uri)
}

// Delete deletes an item of specific uri.
func Delete(uri string) {
	globalBox.Delete(uri)
}

// CheckAndStore checks resp and then store it of specific uri.
func CheckAndStore(uri string, resp *http.Response) {
	globalBox.CheckAndStore(uri, resp)
}

// Clear refresh cache box in d.
func Clear(d time.Duration) {
	globalBox.Clear(d)
}
