FROM golang:1.13 AS builder

RUN mkdir -p /go/src/github.com/cwr0401/redis_metrics

COPY ./ /go/src/github.com/cwr0401/redis_metrics

RUN go get github.com/cwr0401/redis_metrics/cmd/...

FROM alpine:latest

COPY --from=builder /go/bin/redis_metrics   /bin/redis_metrics

ENTRYPOINT ["/bin/redis_metrics"]

