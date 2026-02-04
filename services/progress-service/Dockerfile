FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o progress-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/progress-service .
EXPOSE 8080
CMD ["./progress-service"]
