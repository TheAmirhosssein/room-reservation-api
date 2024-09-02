FROM golang:1.23 AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.23

WORKDIR /app

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /app .

