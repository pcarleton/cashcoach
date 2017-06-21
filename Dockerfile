#FROM golang:onbuild
FROM golang

ADD . /go/src/github.com/pcarleton/cashcoach
ADD ./api/config.json /go/config.json

RUN go get github.com/pcarleton/cashcoach/...

RUN go install github.com/pcarleton/cashcoach/api

ENTRYPOINT /go/bin/api

EXPOSE 5001
