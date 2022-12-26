FROM golang:1.19.3

WORKDIR /otel-simple-page

ADD . /otel-simple-page

RUN cd /otel-simple-page

RUN set GO111MODULE=on

RUN go build -o otel-demo main.go

EXPOSE 8888

ENTRYPOINT ./otel-demo
