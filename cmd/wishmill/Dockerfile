# syntax=docker/dockerfile:1
ARG ARCH="amd64"
ARG OS="linux"
FROM golang:1.19-alpine as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build ./cmd/wishmill

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/wishmill .
RUN apk --no-cache add tzdata
ENTRYPOINT [ "/app/wishmill" ]