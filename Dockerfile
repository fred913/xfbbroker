FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main cmd/main/hello.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
