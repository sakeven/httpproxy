FROM google/golang
MAINTAINER Sakeven "sakeven.jiang@daocloud.io"

# Build app
ENV GOPATH /gopath/app
ADD . /gopath/app/src/httpproxy
WORKDIR /gopath/app/src/httpproxy


RUN go get -t httpproxy
RUN go build httpproxy

EXPOSE 8080
CMD ["/gopath/app/src/httpproxy/httpproxy"]