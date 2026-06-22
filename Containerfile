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

FROM node:24-alpine

RUN apk add --no-cache ca-certificates tzdata supervisor
ENV TZ=Asia/Shanghai

COPY --from=server-builder /order-server /usr/local/bin/order-server
COPY scripts/migrations /app/migrations

COPY --from=web-builder /build/build /app/web/build
COPY --from=web-builder /build/package.json /app/web/

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /app
EXPOSE 3000 8080

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
