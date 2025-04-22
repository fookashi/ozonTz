FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY app/go.mod app/go.sum ./

RUN go mod download

COPY app/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o ozon-app ./cmd/main.go

FROM alpine:latest as runner

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/ozon-app .

CMD ["./ozon-app"]