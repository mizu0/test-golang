FROM golang:1.15-alpine AS builder
WORKDIR /go/src/github.com/mizu0/test-image/
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .


FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/mizu0/test-image/server .
CMD ["./server"]