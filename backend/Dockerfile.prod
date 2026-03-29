FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init && \
    CGO_ENABLED=0 GOOS=linux go build -o backend main.go && \
    chmod +x healthcheck.sh

FROM alpine:3.21
RUN apk add --no-cache ca-certificates lsof
WORKDIR /app

COPY --from=builder /app/backend .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "./healthcheck.sh" ]
CMD ["./backend"]
