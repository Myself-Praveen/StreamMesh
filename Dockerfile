FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o streammesh ./cmd/streammesh

# Final stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/streammesh .

EXPOSE 8080

CMD ["./streammesh"]
