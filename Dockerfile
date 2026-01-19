# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ledger-api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/ledger-api .
# Copy migrations if you plan to run them from the app or a separate tool in the container
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./ledger-api"]
