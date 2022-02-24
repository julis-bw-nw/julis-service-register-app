package data

import (
	"fmt"
	"net"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
)

//go:embed schema.sql
var schema string

type Service struct {
	Host     string
	Database string
	Username string
	Password string
	Timeout  time.Duration

	*pgxpool.Pool
}

func (s Service) dsn() string {
	addr, port, _ := net.SplitHostPort(s.Host)
	return fmt.Sprintf(
		"host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode=disable",
		addr,
		port,
		s.Username,
		s.Password,
		s.Database,
	)
}

func (s *Service) Open() error {
	cfg, err := pgxpool.ParseConfig(s.dsn())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	s.Pool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return err
	}

	return s.Ping(ctx)
}

func (s *Service) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	_, err := s.Exec(ctx, schema)
	return err
}
