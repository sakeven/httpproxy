FROM golang:1.8
MAINTAINER Sakeven Jiang "sakeven.jiang@gmail.com"

# Build app
COPY . $GOPATH/src/github.com/sakeven/httpproxy

RUN go install github.com/sakeven/httpproxy
EXPOSE 8080

CMD ["httpproxy"]
