FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main cmd/main/hello.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main /usr/local/bin

CMD ["/usr/local/bin/main"]
