FROM golang:1.19-alpine3.16 AS builder

# Create and change to the /app directory
WORKDIR /build

# Copy everything to the /build directory
COPY ./ /build

# Download dependencies
#COPY go.mod .
#COPY go.sum .
RUN go mod download

# Build sportsnews
RUN go build -o sportsnews cmd/main.go

# Deploy
FROM alpine:latest

WORKDIR /build

# Copy sportsnews from builder
COPY --from=builder /build/sportsnews .

ENTRYPOINT ["./sportsnews"]