FROM golang

ADD . /go/src/github.com/pcarleton/cashcoach/api
ADD ./config.json /go/config.json

RUN go install github.com/pcarleton/cashcoach/api

ENTRYPOINT /go/bin/api

EXPOSE 5001
