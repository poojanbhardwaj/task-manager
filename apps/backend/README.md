# Task Manager backend

Go API for Task Manager, maintained by poojanbhardwaj.

Module: `github.com/poojanbhardwaj/task-manager/apps/backend`.

## Run

Follow the root README to configure PostgreSQL, Redis, integrations, and database migrations. From this directory:

```powershell
Copy-Item .env.sample .env
# Edit .env before starting.
go run ./cmd/task-manager
```

Configuration keys start with `TASK_MANAGER_` and use dots for nested fields, such as `TASK_MANAGER_DATABASE.HOST`.

## Structure

- `cmd/task-manager`: application entry point.
- `internal/config`: environment configuration.
- `internal/handler`: HTTP handlers.
- `internal/service`: application services.
- `internal/repository`: data access.
- `internal/model`: domain models.
- `internal/database`: PostgreSQL and SQL migrations.
- `internal/middleware`: HTTP middleware.
- `internal/lib`: jobs, email delivery, and utilities.
- `static`: API documentation.
- `templates/emails`: exported email templates.

## Checks

```powershell
go test ./...
go build ./...
```

The starter currently has no Go test files.

## License

The original MIT notices are in the root `LICENSE` file.
