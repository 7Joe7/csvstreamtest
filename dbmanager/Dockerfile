FROM golang:1.11 as builder

WORKDIR /go/src/github.com/7joe7/csvstreamtest/dbmanager
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netcgo -installsuffix cgo .
ENTRYPOINT ["/go/src/github.com/7joe7/csvstreamtest/dbmanager/dbmanager"]

FROM alpine:3.6

RUN until apk --update add ca-certificates; do sleep 1; done
COPY --from=builder /go/src/github.com/7joe7/csvstreamtest/dbmanager/dbmanager /go/bin/dbmanager
ENTRYPOINT ["/go/bin/dbmanager"]