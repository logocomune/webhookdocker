FROM --platform=$BUILDPLATFORM golang:1.23.10-alpine3.22 AS builder

RUN apk add --no-cache git curl build-base  bash


RUN mkdir -p /app

WORKDIR /app

WORKDIR /app

COPY . .


## Set env for multi arch build
ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN echo "GOOS: $GOOS and GOARCH: $GOARCH"


RUN ./build.sh webhook-docker





FROM alpine:3.22

WORKDIR /

COPY --from=builder /app/webhook-docker /application


ENTRYPOINT ["./application"]

