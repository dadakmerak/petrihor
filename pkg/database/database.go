package database

import (
	"fmt"
	"time"

	"github.com/dadakmerak/petrihor/pkg/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Conn *sqlx.DB
}

func (d *DB) InitDatabase(c *config.DB) *DB {
	var err error
	source := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.DBName)

	if d.Conn, err = sqlx.Connect("pgx", source); err != nil {
		panic(err)
	}

	d.Conn.SetConnMaxLifetime(time.Minute * 2)
	d.Conn.SetMaxIdleConns(0)
	d.Conn.SetMaxOpenConns(100)

	return d
}
