package proxy_test

import (
	"httpProxy/proxy"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var srv = httptest.NewServer(nil)

func init() {
	http.DefaultServeMux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, client"))
	}))
}

func ProxyClient() (*http.Client, *httptest.Server) {
	proxy := proxy.NewProxyServer()
	s := httptest.NewServer(proxy)

	proxyUrl, _ := url.Parse(s.URL)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	return client, s
}

func TestProxyUrl(t *testing.T) {
	client, s := ProxyClient()
	defer s.Close()

	txt, err := getTxt(client, srv.URL+"/123")
	if err != nil {
		t.Error(err)
	}
	if txt != "404 page not found\n" {
		t.Error("Getting url error")
	}

	txt, err = getTxt(client, srv.URL+"/hello")
	if err != nil {
		t.Error(err)
	}
	if txt != "Hello, client" {
		t.Error("Getting url error")
	}

}

func getTxt(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(txt), nil
}
