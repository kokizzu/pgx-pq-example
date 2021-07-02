package pgx_pq_example

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func TestPq(t *testing.T) {
	start := time.Now()
	defer func() {
		fmt.Println(`PQ`, time.Since(start))
	}()
	err := godotenv.Load()
	L.PanicIf(err, `failed load .env`)

	pgCfg := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		`127.0.0.1`,
		5432,
		os.Getenv(`POSTGRES_USER`),
		os.Getenv(`POSTGRES_PASSWORD`),
		os.Getenv(`POSTGRES_DB`),
	)

	conn, err := sql.Open(`postgres`, pgCfg)
	conn.SetMaxOpenConns(1)
	L.PanicIf(err, `cannot connect db`)
	defer conn.Close()
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS bar1(id BIGSERIAL PRIMARY KEY, foo VARCHAR(10))`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.Exec(`TRUNCATE TABLE bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	wg := sync.WaitGroup{}
	const K = 1000

	t.Run(`insert`, func(t *testing.T) {
		for z := 0; z < K; z++ {
			wg.Add(1)
			go func(z int) {
				_, err = conn.Exec(`INSERT INTO bar1(foo)VALUES($1)`, I.ToStr(z))
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
				_, err = conn.Exec(`UPDATE bar1 SET foo=$1 WHERE id=$2`, I.ToStr(100-z), z)
				L.PanicIf(err, `failed update bar1`)
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	row := conn.QueryRow(`SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, K, count)
	fmt.Println(count)
}
