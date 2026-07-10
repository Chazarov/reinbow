FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache curl
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]