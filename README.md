

# PGX and PQ Example

Small benchmark example comparing concurrent inserts and updates through
`pgx` and `lib/pq`.

## Requirements

- Go 1.26 or newer
- Docker Compose and PostgreSQL for the integration benchmark

## Usage

Run the default sanity checks without starting PostgreSQL:

```sh
make test
```

Run the integration benchmark:

```sh
docker compose up -d
RUN_POSTGRES_INTEGRATION=1 make integration-test
```

Run a dependency vulnerability check:

```sh
make vulncheck
```

Run any one-off command through the Makefile:

```sh
make run CMD='go test ./...'
```

## Maintenance Checklist

- [x] Update the Go runtime directive to Go 1.26.
- [x] Update the PGX example to `github.com/jackc/pgx/v5`.
- [x] Gate PostgreSQL integration tests behind `RUN_POSTGRES_INTEGRATION=1`.
- [x] Add Makefile targets for tests, integration tests, vulnerability checks, and arbitrary commands.
- [x] Fix update benchmark ids to match PostgreSQL `BIGSERIAL` ids.
