FROM golang:1.27-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

FROM alpine:3.24

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server /app/server

USER appuser:appgroup

EXPOSE 8080
ENTRYPOINT ["/app/server"]
