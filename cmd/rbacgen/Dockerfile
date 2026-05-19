FROM golang:1.26.2-alpine AS builder

ARG ARCH=amd64
ARG GO111MODULE=on

RUN apk add --no-cache git ca-certificates

WORKDIR /code

COPY go.mod go.sum ./
COPY . /code

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build \
    -o /go/bin/rbacgen \
    -trimpath \
    ./cmd/rbacgen

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/rbacgen /bin/rbacgen

ENTRYPOINT ["/bin/rbacgen"]
