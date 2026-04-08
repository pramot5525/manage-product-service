# PRM Product

Simple Go + Fiber API with PostgreSQL and Docker Compose.

## Requirements

- Docker
- Docker Compose
- Go (optional, for local test)

## Quick Start

1. Copy env file:

   cp .env.example .env

2. Start services:

   make start

3. Stop services:

   make down

## Useful Commands

- `make start` : run app + postgres with build
- `make down` : stop and remove containers
- `make restart` : restart all services
- `make logs` : tail app logs
- `make ps` : show compose status
- `make build` : rebuild images without cache
- `make test` : run go tests

## API Docs

- Swagger UI: http://localhost:8080/api-docs/
- OpenAPI file: docs/openapi.yaml
