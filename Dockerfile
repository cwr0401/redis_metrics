FROM golang:1.13-alpine AS builder

RUN mkdir -p /go/src/github.com/cwr0401/redis_metrics

COPY ./ /go/src/github.com/cwr0401/redis_metrics

RUN go get github.com/cwr0401/redis_metrics/cmd/...

FROM alpine:3.10

COPY --from=builder /go/bin/redis-metrics   /usr/local/bin/redis-metrics

ENTRYPOINT ["redis-metrics"]
