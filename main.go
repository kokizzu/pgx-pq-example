package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	L.PanicIf(err, `failed load .env`)
	bg := context.Background() // shared because this not a service
	pgUrl := "postgres://%s:%s@localhost:5432/%s"
	pgUrl = fmt.Sprintf(pgUrl,
		os.Getenv(`POSTGRES_USER`),
		os.Getenv(`POSTGRES_PASSWORD`),
		os.Getenv(`POSTGRES_DB`),
	)

	conn, err := pgx.Connect(bg, pgUrl)
	L.PanicIf(err, `cannot connect db`)
	defer conn.Close(bg)

	_, err = conn.Exec(bg, `CREATE TABLE IF NOT EXISTS bar1(id BIGSERIAL PRIMARY KEY, foo VARCHAR(10))`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.Exec(bg, `TRUNCATE TABLE bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	for z := 0; z < 1000; z++ {
		_, err = conn.Exec(bg, `INSERT INTO bar1(foo)VALUES($1)`, I.ToStr(z))
		L.PanicIf(err, `failed insert to bar1`)
	}

	row := conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	fmt.Println(count)
}
