# Local file storage

Use SeaweedFS instead of a paid Sevalla bucket while following the tutorial.
It exposes an S3-compatible API, so the Go backend still uses the AWS SDK.
PostgreSQL holds task records; SeaweedFS holds attachment contents. Clerk still
handles authentication. You do not need an AWS account or a cloud subscription.

## Windows setup

From the project root:

```powershell
bun storage:setup
bun storage:start
```

Leave that terminal running. In another terminal:

```powershell
bun storage:check
```

The check uploads a uniquely named tiny file, downloads it, tests a temporary
download URL, checks that unsigned access is denied, and deletes the test file.
It does not need PostgreSQL, Redis, or the API server to be running.

The Windows x64 installer downloads SeaweedFS 4.48 from the official GitHub
release and checks its SHA-256 before extracting it. It generates credentials
once and appends only missing storage settings to `apps/backend/.env`.
Rerunning setup preserves your existing settings and stored files.

## Settings and files

| Setting | Local value |
| --- | --- |
| `TASK_MANAGER_STORAGE.ENDPOINT` | `http://127.0.0.1:8333` |
| `TASK_MANAGER_STORAGE.REGION` | `us-east-1` (SDK signing region) |
| `TASK_MANAGER_STORAGE.BUCKET` | `task-manager` |
| `TASK_MANAGER_STORAGE.ACCESS_KEY` | Generated in backend `.env` |
| `TASK_MANAGER_STORAGE.SECRET_KEY` | Generated in backend `.env` |
| `TASK_MANAGER_STORAGE.USE_PATH_STYLE` | `true` |

Path-style addressing puts the bucket in the URL path rather than trying to
resolve a hostname such as `task-manager.localhost`.

The executable is in `.local/seaweedfs/4.48/weed.exe`; uploaded data is in
`.local/storage-data/`. Both are ignored by Git. `bun clean` does not delete
them. Stopping SeaweedFS preserves files; deleting its data directory does not.
Keep these credentials on the backend, never in a `VITE_` variable.

All SeaweedFS listeners are bound to 127.0.0.1. This is a local development
setup, accessible from your computer, and must be running when testing uploads.
SeaweedFS also uses local master/filer/volume ports such as 9333, 8888, and 9340.
If a port is in use, stop your previous storage process before starting another.
Ctrl+C in its terminal stops storage. Restart with `bun storage:start`.

## Backend usage

The API server exposes `server.Storage` when storage settings are provided.
The client provides `Put`, `Get`, `Delete`, and `DownloadURL` (five-minute signed
URLs). `Get` returns a stream; close it after reading.

Storage setup does not add task attachment routes or upload buttons. When you
reach those lessons, check the signed-in user's ownership of a task before
uploading or serving its attachments. Store the object key with the task record.
Never add a public endpoint that signs arbitrary object keys supplied by users.

The tutorial calls its configuration AWS; this project uses `Storage` and
`TASK_MANAGER_STORAGE.*`. Map its region, bucket, endpoint, access key, and
secret key to the fields above. The implementation uses AWS SDK v2
`BaseEndpoint` and `UsePathStyle` to support the local service.

Existing deployments without any storage settings can continue without the
storage client. To later use hosted S3 storage, change the backend settings to
that provider's values, create the bucket there, and configure its CORS and
private access policy. Do not run the local start command for a hosted endpoint.

Other platforms can run SeaweedFS with the same credentials and bucket using
the official quick start; this repository's automatic binary installer is for
Windows x64 only.

References:
- https://github.com/seaweedfs/seaweedfs#quick-start
- https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-endpoints.html
