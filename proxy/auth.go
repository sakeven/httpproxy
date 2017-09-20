package proxy

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// HTTP407 http 407 response
var HTTP407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Secure Proxys\"\r\n\r\n")

// Auth provides basic authorization for proxy server.
func (proxy *Server) Auth(rw http.ResponseWriter, req *http.Request) bool {
	var err error
	if cnfg.Reverse == false && cnfg.Auth == true { //代理服务器登入认证
		if proxy.User, err = proxy.auth(rw, req); err != nil {
			log.Debugf("%v", err)
			return true
		}
	}

	proxy.User = "Anonymous"

	return false
}

func (proxy *Server) auth(rw http.ResponseWriter, req *http.Request) (string, error) {

	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)

	if auth == "" {
		NeedAuth(rw, HTTP407)
		return "", errors.New("need proxy authorization")
	}

	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		log.Debug("when decoding %v, got an error of %v", auth, err)
		return "", errors.New("Fail to decoding Proxy-Authorization")
	}

	var user, passwd string

	userPasswdPair := strings.Split(string(data), ":")
	if len(userPasswdPair) != 2 {
		NeedAuth(rw, HTTP407)
		return "", errors.New("Fail to log in")
	}

	user = userPasswdPair[0]
	passwd = userPasswdPair[1]

	if Check(user, passwd) == false {
		NeedAuth(rw, HTTP407)
		return "", errors.New("Fail to log in")
	}
	return user, nil
}

// NeedAuth requires authorization
func NeedAuth(rw http.ResponseWriter, challenge []byte) error {
	hj, _ := rw.(http.Hijacker)
	client, _, err := hj.Hijack()
	if err != nil {
		return errors.New("Fail to get Tcp connection of Client")
	}
	defer client.Close()

	client.Write(challenge)
	return nil
}

// Check checks username and password
func Check(user, passwd string) bool {
	if user != "" && passwd != "" && cnfg.User[user] == passwd {
		return true
	}
	return false
}
