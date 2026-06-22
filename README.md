# Order - Coffee Shop Ordering System

A modern coffee shop ordering system built with Go and Svelte.

## Tech Stack

- **Backend**: Go + go-chi + protobuf + grpc-gateway + GORM
- **Frontend**: Svelte + SvelteKit + shadcn-svelte
- **Database**: PostgreSQL
- **Auth**: JWT
- **Container**: Podman + Buildah

## Project Structure

```
order/
├── api/order/v1/          # Protobuf definitions
├── cmd/server/            # Application entry point
├── configs/               # Configuration files
├── internal/              # Private application code
│   ├── config/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   └── service/
├── migrations/            # Database migrations
├── pkg/                   # Public library code
│   ├── auth/
│   ├── database/
│   ├── logger/
│   └── response/
├── tools/                 # Custom protoc plugins
├── web/                   # Frontend application
├── Containerfile          # All-in-one container
├── Containerfile.order    # Backend container
├── Containerfile.order-web # Frontend container
├── compose.yml            # Container orchestration
└── Makefile               # Build commands
```

## Getting Started

### Prerequisites

- Go 1.26+
- Node.js 24+
- pnpm
- PostgreSQL 16+
- Podman (optional)

### Development

```bash
# Initialize project
make init

# Start development
make dev

# Build
make build

# Run tests
make test
```

### Docker

```bash
# Build all containers
make docker-build

# Run with all-in-one container
make docker-up

# Stop containers
make docker-down
```

## API Documentation

API is defined using Protocol Buffers. See `api/order/v1/order.proto` for the full definition.

### Endpoints

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `GET /api/v1/products` - List products
- `GET /api/v1/products/:id` - Get product
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List user orders
- `GET /api/v1/orders/:id` - Get order details

## License

MIT
