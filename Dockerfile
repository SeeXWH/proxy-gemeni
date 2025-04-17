# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o gemini-app

# Создаем итоговый образ
FROM alpine:latest

# Устанавливаем зависимости для runtime
RUN apk --no-cache add ca-certificates

# Копируем бинарный файл из builder
WORKDIR /root/
COPY --from=builder /app/gemini-app .

# Открываем порт, который будет использовать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./gemini-app"]