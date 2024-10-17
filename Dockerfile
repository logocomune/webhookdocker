FROM golang:1.23.2-alpine3.20 AS builder

RUN apk add --no-cache git curl build-base  bash


RUN mkdir -p /app

WORKDIR /app

WORKDIR /app

COPY . .

RUN ./build.sh webhook-docker





FROM alpine:3.20

WORKDIR /

COPY --from=builder /app/webhook-docker /application


ENTRYPOINT ["./application"]

