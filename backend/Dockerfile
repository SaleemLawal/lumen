# Hot reload for local Compose (build.target: dev)
FROM golang:1.25.1-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git ca-certificates make

RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download && go mod verify

CMD ["air"]

FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./cmd/api

FROM alpine:latest AS production

COPY --from=builder /app/main /usr/local/bin/main

CMD ["/usr/local/bin/main"]
