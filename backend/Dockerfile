FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    swag init -g internal/cmd/main.go && \
    go build -o backend internal/cmd/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates lsof
WORKDIR /app

ENV GIN_MODE=release

COPY --from=builder /app/backend .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./backend"]
