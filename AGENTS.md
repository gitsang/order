# AGENTS.md - Order Coffee Shop

## Critical Build Rule

**NEVER compile binaries to the current directory.** They will be committed accidentally.

```bash
# Test compilation (preferred)
go build ./...

# If binary is needed, output to safe locations
go build -o ./bin/server cmd/server/main.go
go build -o /tmp/server cmd/server/main.go
```

The Makefile uses `bin/` which is in `.gitignore`.

## Project Structure

```
order/
├── cmd/server/main.go      # Go entrypoint
├── internal/               # Private Go packages
│   ├── config/             # Viper config (searches ./configs/)
│   ├── handler/            # HTTP handlers + middleware
│   ├── model/              # GORM models
│   ├── repository/         # Database operations
│   └── service/            # Business logic
├── configs/                # YAML config (gitignored except examples)
├── scripts/migrations/     # SQL migrations (golang-migrate)
├── web/                    # SvelteKit frontend
│   └── src/
│       ├── lib/api/        # API client layer
│       ├── lib/stores/     # Svelte writable stores
│       └── routes/         # File-based routing
└── compose.yml             # PostgreSQL + app containers
```

## Developer Commands

```bash
# Init (first time)
make init                   # go mod tidy + pnpm install + start postgres

# Development
make dev                    # Run Go server (port 8080)
make dev-web                # Run Svelte dev server (port 5173)

# Build
make build                  # Build both (outputs to bin/ and web/build/)
make build-server           # Go → bin/server
make build-web              # Svelte → web/build/

# Quality
make lint                   # golangci-lint + pnpm lint
make test                   # go test ./... + pnpm test

# Database migrations
make migrate-up DB_URL="postgres://order:order123@localhost:5432/order?sslmode=disable"
make migrate-down DB_URL="..." 
make migrate-create NAME=add_users_table
```

## Architecture Quirks

### Config System
- Viper loads `config.yaml` from `./`, `./configs/`, or `/etc/order/`
- Environment variables override file values (e.g., `DB_HOST`, `JWT_SECRET`)
- Default DB credentials: `order` / `order123`

### Auth & Roles
- JWT tokens contain `user_id`, `username`, `role` claims
- Roles: `customer` (default), `admin`
- `AdminMiddleware` in `internal/handler/middleware.go` gates admin routes
- Frontend checks `$auth.user?.role === 'admin'` for admin access

### Frontend (SvelteKit)
- **Svelte 5 runes mode** enabled globally (`.svelte` files use `$state`, `$derived`, `$effect`)
- `adapter-node` serves production build via `node web/build`
- API client at `web/src/lib/api/client.ts` auto-attaches JWT from localStorage
- 401 responses trigger redirect to `/login`
- Stores use `browser` check to avoid SSR localStorage errors

### Database
- PostgreSQL with UUID primary keys (`gen_random_uuid()`)
- GORM soft deletes (`deleted_at` column)
- Models in `internal/model/model.go` have `BeforeCreate` hooks

### Deployment
- All-in-one container uses supervisord to run both `node /app/web/build` (port 3000) and `order-server` (port 8080)
- Split containers available via `--profile split`

## API Response Format

All endpoints return:
```json
{
  "code": 0,        // 0 = success, error code otherwise
  "message": "ok",
  "data": { ... }
}
```

Helper: `response.Success(w, data)` and `response.Error(w, statusCode, message)`

## File Conventions

- Config examples: `configs/config.example.yml` (tracked)
- Actual configs: `configs/config.yml` (gitignored)
- Migrations: `scripts/migrations/NNN_description.sql`
- Go tests: colocated `*_test.go` files
- Frontend types: `web/src/lib/api/types.ts` is single source of truth

## Common Pitfalls

1. **Don't use `go build` without `-o` flag** — binary lands in repo root
2. **Config not found** — check you're running from project root or `configs/` exists
3. **Frontend 401 loops** — `client.ts` redirects to `/login` on 401, ensure login page exists
4. **SSR localStorage errors** — always wrap in `if (browser)` check
5. **Svelte 5 syntax** — use `$state()` not `let x = writable()`, use `onclick` not `on:click`
