FROM golang:1.26.2-alpine AS builder

# Set environment variables for Go module and x64 architecture
ARG ARCH=amd64
ARG GO111MODULE=on

RUN apk add --no-cache git ca-certificates

WORKDIR /code

# Copy go mod and sum files and download all the dependencies
COPY go.mod go.sum ./
COPY . /code

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build \
    -o /go/bin/user-management \
    -trimpath \
    ./cmd/user-management

# Run the container
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /go/bin/user-management /bin/user-management

# Copy config file
COPY --from=builder /code/cmd/user-management/config.yaml ./config.yaml

EXPOSE 3301

ENTRYPOINT ["/bin/user-management"]