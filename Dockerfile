FROM golang:1.15.7 as builder

WORKDIR /infura-exporter
COPY . .

ARG VERSION=undefined
RUN CGO_ENABLED=0 \
    go build

FROM alpine:latest
ENTRYPOINT ["/infura-exporter"]
USER nobody
EXPOSE 9877

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /infura-exporter/infura-exporter /infura-exporter
COPY config-example.yaml /config.yaml