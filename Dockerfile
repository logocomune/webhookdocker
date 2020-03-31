FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git curl make binutils

RUN mkdir -p /app

WORKDIR /app

ENV GOPROXY="https://proxy.golang.org,direct"

COPY go.* ./

RUN go mod download


FROM build_base as builder

COPY . .

RUN make docker_install

RUN strip /go/bin/webhook-docker


FROM alpine

WORKDIR /

COPY --from=builder /go/bin/webhook-docker /webhook-docker

ENTRYPOINT ["./webhook-docker"]

