package config

import (
	"context"
	"fmt"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func NewPostgres(ctx context.Context, cfg Postgres) *pgxpool.Pool {
	pool, err := pgxpool.Connect(ctx, cfg.String())
	fmt.Println("config database: ", cfg.String())
	if err != nil {
		panic("no connect to database")
	}

	return pool
}

func (p *Postgres) String() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:   p.DBName,
	}

	q := u.Query()
	q.Set("sslmode", "disable")

	u.RawQuery = q.Encode()

	return u.String()
}

func (p *Postgres) MigrationsUp(url ...string) error {
	var sourceURL string
	if url == nil {
		sourceURL = "file://internal/migrations/up"
	} else {
		sourceURL = url[0]
	}
	m, err := migrate.New(sourceURL, p.String())
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		return err
	}

	return nil
}
