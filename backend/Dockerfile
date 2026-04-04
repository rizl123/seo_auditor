# --- Этап сборки ---
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Собираем бинарник (swag init больше не нужен)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o backend internal/cmd/main.go

# --- Финальный образ ---
FROM alpine:3.21
RUN apk add --no-cache ca-certificates lsof
WORKDIR /app

# Копируем только исполняемый файл
COPY --from=builder /app/backend .

# ENV GIN_MODE нам больше не нужен, так как мы на стандартном http mux
EXPOSE 8080

CMD ["./backend"]
