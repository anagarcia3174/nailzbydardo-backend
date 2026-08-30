# --- Build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependency downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# --- Run stage ---
FROM alpine:latest

WORKDIR /app

# Needed for HTTPS calls out of the container (e.g. if you ever add
# outbound requests) and for TLS verification generally
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]