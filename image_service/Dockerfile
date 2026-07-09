# syntax=docker/dockerfile:1

FROM golang:1.26.4-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    pkg-config \
    libmagickwand-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /imageservice ./cmd/image_service/

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    imagemagick \
    fonts-dejavu-core \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /imageservice .

RUN mkdir -p /data/images

EXPOSE 8080

CMD ["./imageservice"]
