package pg

import (
	"MegaCode/internal/pkg/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	conn pgx.Conn
}

func NewDataBase(cfg PgConfig) (Database, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/postgres?sslmode=disable", cfg.name, cfg.password, cfg.port)
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return Database{}, err
	} else {
		return Database{*conn}, nil
	}
}
func (d *Database) Insert(userinfo model.User) error {
	_, err := d.conn.Exec(context.Background(), "insert to users values ($1, $2)", userinfo.Email, userinfo)
	return err
}

func (d *Database) Select(userEmail string) (string, error) {
	rows, err := d.conn.Exec(context.Background(), "SELECT password FROM users WHERE email=($1)", userEmail)
	if err != nil {
		return "", err
	}
	return rows.String(), err
}
