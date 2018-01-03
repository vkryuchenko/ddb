/*
author Vyacheslav Kryuchenko
*/
package helpers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	SCHEMA = `
CREATE TABLE IF NOT EXISTS users
(
  id     SERIAL                NOT NULL
    CONSTRAINT users_pkey
    PRIMARY KEY,
  login  VARCHAR(32)           NOT NULL,
  email  VARCHAR(64)           NOT NULL,
  active BOOLEAN DEFAULT TRUE  NOT NULL,
  admin  BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
  ON users (login);

CREATE TABLE IF NOT EXISTS history
(
  action_time TIMESTAMP DEFAULT now() NOT NULL
    CONSTRAINT history_pkey
    PRIMARY KEY,
  user_id     INTEGER                 NOT NULL
    CONSTRAINT history_users__fk
    REFERENCES users
    ON UPDATE CASCADE ON DELETE CASCADE,
  command     TEXT                    NOT NULL,
  success     BOOLEAN DEFAULT TRUE    NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
  user_id INTEGER                                            NOT NULL
    CONSTRAINT sessions_pkey
    PRIMARY KEY
    CONSTRAINT sessions_users__fk
    REFERENCES users (id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  expire  TIMESTAMP DEFAULT (now() + '08:00:00' :: INTERVAL) NOT NULL
);
`
)

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

func (pg *PostgresConfig) InitDatabase() error {
	sqlCon, err := pg.connect()
	if err != nil {
		return err
	}
	_, err = sqlCon.Exec(SCHEMA)
	if err != nil {
		return err
	}
	return nil
}
