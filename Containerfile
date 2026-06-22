FROM golang:1.26-alpine AS server-builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /order-server cmd/server/main.go

FROM node:24-alpine AS web-builder

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /build
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ .
RUN pnpm build

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata nginx supervisor
ENV TZ=Asia/Shanghai

COPY --from=server-builder /order-server /usr/local/bin/order-server
COPY migrations /app/migrations

COPY --from=web-builder /build/build /var/www/html

COPY nginx.conf /etc/nginx/http.d/default.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /app
EXPOSE 80 8080

CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
