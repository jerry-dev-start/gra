# ---- 构建阶段 ----
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gra ./cmd/server

# ---- 运行阶段 ----
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/gra .
COPY --from=builder /app/config/config-prod.yaml ./config/config.yaml

EXPOSE 8888

CMD ["./gra"]
