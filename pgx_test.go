package pgx_pq_example

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"

	"github.com/joho/godotenv"
)

func TestPgx(t *testing.T) {
	start := time.Now()
	defer func() {
		fmt.Println(`PGX`, time.Since(start))
	}()
	err := godotenv.Load()
	L.PanicIf(err, `failed load .env`)

	bg := context.Background() // shared context
	pgUrl := "postgres://%s:%s@%s:%d/%s"
	pgUrl = fmt.Sprintf(pgUrl,
		os.Getenv(`POSTGRES_USER`),
		os.Getenv(`POSTGRES_PASSWORD`),
		`127.0.0.1`,
		5432,
		os.Getenv(`POSTGRES_DB`),
	)

	conn, err := pgxpool.Connect(bg, pgUrl)
	L.PanicIf(err, `cannot connect db`)
	defer conn.Close()

	_, err = conn.Exec(bg, `CREATE TABLE IF NOT EXISTS bar1(id BIGSERIAL PRIMARY KEY, foo VARCHAR(10))`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.Exec(bg, `TRUNCATE TABLE bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	wg := sync.WaitGroup{}
	const K = 1000

	t.Run(`insert`, func(t *testing.T) {
		for z := 0; z < K; z++ {
			wg.Add(1)
			go func(z int) {
				_, err = conn.Exec(bg, `INSERT INTO bar1(foo)VALUES($1)`, I.ToStr(z))
				L.PanicIf(err, `failed insert to bar1`)
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	t.Run(`update`, func(t *testing.T) {
		for z := 0; z < K; z++ {
			wg.Add(1)
			go func(z int) {
				_, err = conn.Exec(bg, `UPDATE bar1 SET foo=$1 WHERE id=$2`, I.ToStr(100-z), z)
				L.PanicIf(err, `failed update bar1`)
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	row := conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, K, count)
	fmt.Println(count)
}
