FROM golang:1.15 AS builder
WORKDIR /go/src/github.com/mizu0/test-image/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .


FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/mizu0/test-image/server .
CMD ["./server"]