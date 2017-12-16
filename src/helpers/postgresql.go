/*
author Vyacheslav Kryuchenko
*/
package helpers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const PGCREATESHEMA = ``

type PostgresConfig struct {
	Address string `json:"address"`
	Auth    string `json:"auth"`
	DBName  string `json:"dbname"`
	SSL     string `json:"ssl,omitempty"`
}

func (pg *PostgresConfig) connect() (*sql.DB, error) {
	var ssl string
	if pg.SSL != "" {
		ssl = pg.SSL
	} else {
		ssl = "disable"
	}
	connectionParams := fmt.Sprintf("postgres://%s@%s/%s?sslmode=%s", pg.Auth, pg.Address, pg.DBName, ssl)
	return sql.Open("postgres", connectionParams)
}
