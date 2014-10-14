#!bin/bash

echo "test without proxy"
curl -I http://www.baidu.com

go build

mv config/config.json config/config.json.bak
mv config/config2.json config/config.json

./httpproxy&

#set http proxy
export http_proxy="http://proxy:proxy@localhost:8080/"
export https_proxy="http://proxy:proxy@localhost:8080/"


sleep 1

#test http connect
echo 'test gfwlist'
curl -I http://www.baidu.com/
curl -I http://www.baidu.com/img/bdlogo.png

echo 
echo 'test no-gfwlist website'
curl -I http://www.163.com/

#kill process httpproxy
pkill httpproxy

mv config/config.json config/config2.json
mv config/config.json.bak config/config.json

echo 'test succeed'