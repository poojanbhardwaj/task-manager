# Task Manager

Task Manager is a full-stack application maintained by **poojanbhardwaj**, with a Go API and a React/TypeScript frontend.

## Project layout

- `apps/backend`: Go API, database migrations, background jobs, and email delivery.
- `apps/frontend`: React application.
- `packages/zod`: shared schemas (`@poojanbhardwaj/zod`).
- `packages/openapi`: API contracts and documentation (`@poojanbhardwaj/openapi`).
- `packages/emails`: email templates (`@poojanbhardwaj/emails`).

## Deployment

- **Frontend:** [https://frontend-production-f7d0.up.railway.app](https://frontend-production-f7d0.up.railway.app)
- **Backend health:** [https://task-manager-production-dcdd.up.railway.app/status](https://task-manager-production-dcdd.up.railway.app/status)
- **API documentation:** [https://task-manager-production-dcdd.up.railway.app/docs](https://task-manager-production-dcdd.up.railway.app/docs)

The frontend and backend are deployed as separate Railway services. Clerk is currently in Development mode for demo and test use. Clerk Production mode requires a custom domain; Railway-generated `*.up.railway.app` domains are not accepted by Clerk for production.

Backend module: `github.com/poojanbhardwaj/task-manager/apps/backend`.

## Prerequisites

- Go 1.24.5 or later.
- Node.js 22.12 or later and Bun.
- PostgreSQL and Redis to run the API.
- Your own Clerk, Resend, and optional New Relic settings when using those integrations.
- Tern for manual database migrations. Task is optional.

## Install

From the project root:

```powershell
bun install --frozen-lockfile
Set-Location apps/backend
go mod download
Set-Location ../..
Copy-Item apps/backend/.env.sample apps/backend/.env
Copy-Item apps/frontend/.env.example apps/frontend/.env.local
```

Edit the copied configuration files with your own settings.

Backend variables use the `TASK_MANAGER_` prefix and dots for nested fields:

```dotenv
TASK_MANAGER_DATABASE.NAME="task_manager"
TASK_MANAGER_REDIS.ADDRESS="localhost:6379"
TASK_MANAGER_SERVER.PORT="8080"
```

Create the PostgreSQL database `task_manager`. Use a Redis host and port without a URL scheme. Configure the frontend Clerk publishable key before using the authenticated API client; the starter screen itself does not call the API.

In local mode, apply migrations from `apps/backend` after installing Tern:

```powershell
tern migrate -m ./internal/database/migrations --conn-string "postgres://postgres:YOUR_PASSWORD@localhost:5432/task_manager?sslmode=disable"
```

Replace the example connection settings. The API applies migrations automatically in non-local environments.

## Develop

Start PostgreSQL and Redis. In one terminal, from the project root:

```powershell
Set-Location apps/backend
go run ./cmd/task-manager
```

In a second terminal, from the project root:

```powershell
bun run dev
```

The root development command starts the frontend, shared-package watchers, and email preview. Start the Go API separately.

- Frontend: http://localhost:3000
- API health: http://localhost:8080/status
- API reference: http://localhost:8080/docs
- Email preview: http://localhost:3001

## Local attachment storage (no cloud subscription)

Use SeaweedFS to follow the tutorial's S3 file-storage section locally:

```powershell
bun storage:setup
bun storage:start
```

Keep that terminal open and run `bun storage:check` in another terminal.
Setup generates private credentials in `apps/backend/.env`; files persist in
the ignored `.local/storage-data` folder. PostgreSQL still stores task records.
See [local storage instructions](docs/local-storage.md) for settings and usage.
Task attachment endpoints and UI are still to be implemented.

## Build and check

From the project root:

```powershell
bun run build
Set-Location apps/backend
go test ./...
go build ./...
```

The starter currently has no Go test files; `go test ./...` checks package compilation and the default vet checks.

After building shared packages, run `bun run gen` from `packages/openapi` to regenerate API documentation. Run `bun run export` from `packages/emails` to export email templates.

## Clean generated files

From the project root, run `bun clean` to remove frontend/package build outputs,
TypeScript build caches, Turbo caches, and installed `node_modules` folders.
Source files, environment files, lockfiles, and Go module caches are preserved.
Cleanup uses Node's built-in filesystem functions and works without installing
JavaScript dependencies first. Running it again is safe.

Run `bun install --frozen-lockfile` afterward before building or developing.
To clean only one workspace, run `bun clean` from that workspace's directory.

## Repository

The project is hosted in the private [GitHub repository](https://github.com/poojanbhardwaj/task-manager).

## License

The original MIT copyright and license notices are preserved in `LICENSE`. Include applicable notices with redistributed source and builds. Dependencies retain their own license terms.
