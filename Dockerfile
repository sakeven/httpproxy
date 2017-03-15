FROM golang:1.7
MAINTAINER Sakeven Jiang "jc5930@sina.cn"

# Build app
ADD . $GOPATH/src/httpproxy
WORKDIR $GOPATH/src/httpproxy

RUN go get -t httpproxy
RUN go build httpproxy

EXPOSE 8080
CMD ["$GOPATH/src/httpproxy/httpproxy"]
