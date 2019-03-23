FROM golang:1.11-alpine as builder

WORKDIR /go/src/staticweb
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

RUN go get -u github.com/palager/staticweb/cmd/staticweb

FROM alpine:latest
COPY --from=builder /go/bin/staticweb /bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/bin/staticweb"]

