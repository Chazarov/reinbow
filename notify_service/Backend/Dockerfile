# ================== ЭТАП СБОРКИ ==================
FROM golang:1.26.4-alpine AS builder

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/backend ./main.go

# Устанавливаем migrate и postgresql-client (для pg_isready)
RUN apk add --no-cache curl postgresql-client && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

# ================== ФИНАЛЬНЫЙ ОБРАЗ ==================
FROM alpine:latest

WORKDIR /app

# Устанавливаем postgresql-client в финальный образ (чтобы был libpq.so и pg_isready)
RUN apk add --no-cache postgresql-client

# Копируем бинарники
COPY --from=builder /app/backend .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/

# pg_isready уже есть в финальном образе после установки postgresql-client

COPY config ./config
COPY migrations ./migrations

EXPOSE 8080

# CMD с ожиданием БД через pg_isready (теперь работает)
CMD sh -c '. ./config/data.env && \
           echo "DATABASE_CONNECT=$DATABASE_CONNECT" && \
           echo "Waiting for database..." && \
           until pg_isready -h db -U postgres; do sleep 2; done && \
           echo "Database ready, running migrations..." && \
           migrate -source file://./migrations -database "$DATABASE_CONNECT" up && \
           echo "Migrations done, starting backend..." && \
           exec ./backend'