FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Устанавливаем git для зависимостейы
RUN apk add --no-cache git

# Копируем весь код сначала
COPY . .

# Если go.mod и go.sum не существуют - создаем, иначе используем существующие
RUN if [ ! -f go.mod ]; then \
        echo "go.mod not found, initializing..."; \
        go mod init buch; \
    fi && \
    go mod tidy

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ (легкий)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Копируем статику и шаблоны
COPY --from=builder /app/src ./src

EXPOSE 8080

CMD ["./main"]
