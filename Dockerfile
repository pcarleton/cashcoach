#FROM golang:onbuild
FROM golang

ADD . /go/src/github.com/pcarleton/cashcoach
ADD ./app/config.json /go/config.json

RUN go get github.com/pcarleton/cashcoach/...

RUN go install github.com/pcarleton/cashcoach/app

ENTRYPOINT /go/bin/app

EXPOSE 5001
