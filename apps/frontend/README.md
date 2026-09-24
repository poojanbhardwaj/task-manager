# Task Manager frontend

React and TypeScript frontend for Task Manager, maintained by poojanbhardwaj.

From the repository root:

```powershell
bun install --frozen-lockfile
Copy-Item apps/frontend/.env.example apps/frontend/.env.local
bun run dev
```

The development server runs at http://localhost:3000. Start the Go API separately using the root README.

Set your Clerk publishable key before using the authenticated API client. The current starter screen displays Task Manager; task-management features are not implemented yet.

Build all frontend dependencies and the app from the repository root:

```powershell
bun run build
```

Shared imports use `@poojanbhardwaj/openapi` and `@poojanbhardwaj/zod`.
