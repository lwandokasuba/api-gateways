# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/service ./main.go

FROM alpine:3.20
RUN adduser -D -g '' app
USER app
COPY --from=build /out/service /usr/local/bin/service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/service"]
