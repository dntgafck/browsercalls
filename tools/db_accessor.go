package tools

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

var mu = &sync.Mutex{}
var dbAccessor *sql.DB

var connParams = map[string]string{
	"dbname":   os.Getenv("BC_DB_NAME"),
	"host":     os.Getenv("BC_DB_HOST"),
	"port":     os.Getenv("BC_DB_PORT"),
	"user":     os.Getenv("BC_DB_USER"),
	"password": os.Getenv("BC_DB_PASSWORD"),
}

func GetDbAccessor(host string, port int, dbname, user, password string) (*sql.DB, error) {
	defer mu.Unlock()
	mu.Lock()
	if nil != dbAccessor {
		return dbAccessor
	}

	dsn, err := createDsn()
	if nil != err {
		return nil, err
	}

	dbAccessor, err = sql.Open("postgres", dsn)

	return dbAccessor, err
}

func createDsn() (dsn string, err error) {
	for k, v := range connParams {
		if 0 == len(v) {
			errText := fmt.Sprintf("Missing db connection param '%s'", k)
			err = errors.New(errText)
			return
		}
	}

	dsn = fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		connParams["dbname"],
		connParams["user"],
		connParams["password"],
		connParams["host"],
		connParams["port"],
	)
	return
}
