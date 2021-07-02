

# PGX and PQ Example

```
alias time='/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB"'

docker-compose up

go mod tidy

go test -v

=== RUN   TestPgx
=== RUN   TestPgx/insert
=== RUN   TestPgx/update
1000
PGX 857.408967ms
--- PASS: TestPgx (0.86s)
    --- PASS: TestPgx/insert (0.81s)
    --- PASS: TestPgx/update (0.03s)
=== RUN   TestPq
=== RUN   TestPq/insert
=== RUN   TestPq/update
1000
PQ 3.399219416s
--- PASS: TestPq (3.40s)
    --- PASS: TestPq/insert (3.14s)
    --- PASS: TestPq/update (0.24s)
PASS
ok      pgx_pq_example  4.261s

```
