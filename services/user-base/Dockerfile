# build stage
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service ./services/user-base/main/

# final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 50051

CMD ["./service"]