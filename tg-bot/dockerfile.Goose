# Используем легкий образ Go
FROM golang:1.25 AS builder

# Рабочая директория
WORKDIR /app

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Финальный легкий образ
FROM debian:12-slim

# Устанавливаем базовые утилиты
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Копируем бинарник goose из builder
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Рабочая директория
WORKDIR /app

# По умолчанию запускается goose (можно менять через docker-compose)
CMD ["goose"]