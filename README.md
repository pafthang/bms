# BMS Backend

BMS is a Go backend service for multi-workspace bookmark management.
It provides:

- Email/password auth with access/refresh tokens
- User profile and user settings APIs
- Workspace CRUD and membership management with roles
- Bookmark and tag CRUD inside workspaces
- OpenAPI generation from registered `arc` routes

Repository includes:

- Backend API (`cmd/api`)
- Migration CLI (`cmd/migrate`)
- OpenAPI generator CLI (`cmd/arc`)
- Web SDK (`web/sdk`)
- Web UI (`web/ui`)

## Tech Stack

- Go
- `github.com/pafthang/arc` for HTTP runtime and OpenAPI generation
- `github.com/pafthang/orm` + `github.com/pafthang/dbx` for data access
- SQLite as the only database

## Requirements

- Go 1.25+

## Environment Variables

| Name | Default | Description |
|---|---|---|
| `BMS_HTTP_ADDR` | `:8080` | HTTP listen address |
| `BMS_DB_DSN` | `file:bms.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)` | SQLite DSN |
| `BMS_JWT_SECRET` | `dev-secret-change-me` | JWT signing secret |
| `BMS_ACCESS_TTL_SEC` | `900` | Access token TTL in seconds |
| `BMS_REFRESH_TTL_SEC` | `2592000` | Refresh token TTL in seconds |
| `BMS_ENABLE_DOCS` | `true` | Enable `/docs` and `/openapi.json` routes |
| `BMS_CORS_ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:5173` | Comma-separated CORS origins |

## Quick Start

### 1. Apply migrations

```bash
go run ./cmd/migrate -action up -dir migrations
```

### 2. Run API server

```bash
go run ./cmd/api
```

API base path:

- `http://localhost:8080/api/v1`

System routes:

- `http://localhost:8080/health`
- `http://localhost:8080/ready`
- `http://localhost:8080/openapi.json`
- `http://localhost:8080/docs`

### 3. Generate OpenAPI

```bash
go run ./cmd/arc -format json -out openapi/openapi.json
```

Quality-gated generation:

```bash
go run ./cmd/arc -format json -out openapi/openapi.json \
  -validate-quality \
  -require-tags \
  -require-servers \
  -require-security-schemes BearerAuth \
  -require-examples
```

### 4. Generate Web SDK (TypeScript)

```bash
pnpm install
pnpm --filter @pafthang/bms-sdk run generate
```

### 5. Run Web UI (SvelteKit)

```bash
pnpm install
pnpm --filter @pafthang/bms-ui run dev
```

### 6. Default admin seed

Migration `008_seed_default_admin` creates a default superadmin (idempotent):

- Email: `admin@admin.admin`
- Password: `adminadminadmin`

## Migrations

Run migration status:

```bash
go run ./cmd/migrate -action status -dir migrations
```

Validate migration set:

```bash
go run ./cmd/migrate -action validate -dir migrations
```

Rollback one step:

```bash
go run ./cmd/migrate -action down -steps 1 -dir migrations
```

## Tests

Run all tests:

```bash
go test ./...
```

The repository includes integration tests using SQLite and real HTTP handlers.

CI additionally verifies OpenAPI quality gates and OpenAPI 3.1 validation.

## License

MIT. See [`LICENSE`](LICENSE).

## API Groups

Current route groups:

- `Auth`
- `Users`
- `UserSettings`
- `Workspaces`
- `WorkspaceUsers`
- `Bookmarks`
- `Tags`
- `Admin`
