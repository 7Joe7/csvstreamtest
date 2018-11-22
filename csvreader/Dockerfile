FROM golang:1.11 as builder

WORKDIR /go/src/github.com/7joe7/csvstreamtest/csvreader
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netcgo -installsuffix cgo .
ENTRYPOINT ["/go/src/github.com/7joe7/csvstreamtest/csvreader/csvreader"]

FROM alpine:3.6

RUN until apk --update add ca-certificates; do sleep 1; done
COPY --from=builder /go/src/github.com/7joe7/csvstreamtest/csvreader/csvreader /go/bin/csvreader
ENTRYPOINT ["/go/bin/csvreader"]
