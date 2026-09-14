FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy go.mod and go.sum (if it exists)
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o streammesh ./cmd/streammesh

# Minimal runtime image
FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/streammesh .

# Copy config if any
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./streammesh"]
