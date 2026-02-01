FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .
RUN go env -w GO111MODULE=off
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]
