FROM golang:1.9 as builder

WORKDIR /go/src/github.com/SantoDE/varaco

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o dist/varaco .

FROM scratch
COPY --from="builder" /go/src/github.com/SantoDE/varaco/dist/varaco /varaco
COPY script/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/varaco"]