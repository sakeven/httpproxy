#!bin/bash

go build

./httpproxy&

#set http proxy
export http_proxy="http://proxy:proxy@localhost:8080/"
export https_proxy="http://proxy:proxy@localhost:8080/"

sleep 1

#test http connect
echo 'test http connect'
curl -I 'http://www.baidu.com'
#test https connect
echo 'test https connect'
curl -I 'https://byvoid.com/'
#download picture
echo 'test download'
curl -I http://www.baidu.com/img/bdlogo.png

#kill process httpproxy
pkill httpproxy

echo 'test succeed'