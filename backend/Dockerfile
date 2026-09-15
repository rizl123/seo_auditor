FROM golang:1.26.1-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o backend internal/cmd/main.go

FROM alpine:3.21
RUN --mount=type=cache,target=/var/cache/apk apk add ca-certificates lsof
WORKDIR /app

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
