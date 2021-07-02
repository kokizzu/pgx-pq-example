package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	TestPq()
	TestPgx()
}

func TestPgx() {
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
	for z := 0; z < 1000; z++ {
		wg.Add(1)
		go func(z int) {
			ctx := context.Background()
			_, err = conn.Exec(ctx, `INSERT INTO bar1(foo)VALUES($1)`, I.ToStr(z))
			L.PanicIf(err, `failed insert to bar1`)
			wg.Done()
		}(z)
	}

	wg.Wait()
	row := conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	fmt.Println(count)
}

func TestPq() {
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
	for z := 0; z < 1000; z++ {
		wg.Add(1)
		go func(z int) {
			_, err = conn.Exec(`INSERT INTO bar1(foo)VALUES($1)`, I.ToStr(z))
			L.PanicIf(err, `failed insert to bar1`)
			wg.Done()
		}(z)
	}

	wg.Wait()
	row := conn.QueryRow(`SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	fmt.Println(count)
}
