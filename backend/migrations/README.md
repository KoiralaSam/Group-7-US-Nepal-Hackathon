# Database migrations

This folder holds [golang-migrate](https://github.com/golang-migrate/migrate) SQL migrations. Each version uses paired files:

- `NNNNNN_name.up.sql` — apply the change  
- `NNNNNN_name.down.sql` — roll back that version  

## Prerequisites

1. **PostgreSQL** running locally (or point `DATABASE_URL` at your instance).
2. **migrate CLI** installed. Either:

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

   Ensure `$(go env GOPATH)/bin` is on your `PATH` so the `migrate` command is found.

3. **Database exists** (migrate does not create the database, only the schema inside it):

   ```bash
   createdb g7
   ```

## Configuration

Set `DATABASE_URL` when you run Make. Defaults match a typical local Postgres setup (user/password/database `postgres`/`postgres`/`g7`).

From the **repository root**:

```bash
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/g7?sslmode=disable'
```

Adjust user, password, host, port, database name, and `sslmode` as needed.

## Running migrations (Makefile)

From the **repository root**, with `migrate` on your `PATH`:

```bash
# Apply all pending migrations
make migrate-up

# Roll back the last applied migration
make migrate-down

# Show current migration version (and dirty flag if any)
make migrate-version
```

One-off invocation without Make:

```bash
migrate -path backend/migrations -database "$DATABASE_URL" up
migrate -path backend/migrations -database "$DATABASE_URL" down 1
```

## Adding a new migration

1. Create the next numbered pair, e.g. `000002_add_sessions.up.sql` and `000002_add_sessions.down.sql`.
2. Keep `down` logic able to undo `up` cleanly.
3. Run `make migrate-up` to apply.

## Troubleshooting

- **`dirty` database version**: a previous run failed mid-migration. Inspect the DB, fix manually if needed, then use `migrate force VERSION` (see golang-migrate docs) or restore from backup.
- **`migrate: command not found`**: install the CLI and add Go’s `bin` directory to `PATH`.
- **Connection refused**: confirm Postgres is listening and `DATABASE_URL` matches host/port.
