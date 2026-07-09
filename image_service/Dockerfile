FROM golang:1.26.4-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    pkg-config \
    libmagickwand-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# ENV CGO_ENABLED=1
RUN go mod tidy
RUN go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g spec.go -o api -d api,internal/app/ports/api/v1 --parseInternal --instanceName internal
RUN go build -o /imageservice ./cmd/image_service/

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    imagemagick \
    fonts-dejavu-core \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /imageservice .

RUN mkdir -p /data/images

EXPOSE 8080

CMD ["./imageservice"]
