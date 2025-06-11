FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install make and protoc dependencies
RUN apk add --no-cache make protobuf-dev protoc

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make proto build

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/study-server .

# Set build argument for environment
ARG ENV=prod
ENV APP_ENV=${ENV}

# Copy the appropriate environment file
COPY --from=builder /app/.env.${ENV} .env

# Expose the port Heroku will use
ENV PORT=1973
EXPOSE 1973

# Run the application
CMD ["./study-server"] 