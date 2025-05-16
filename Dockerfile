FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bot-app ./main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=builder /app/bot-app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./

CMD ["./bot-app"]