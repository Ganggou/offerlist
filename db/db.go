package db

import (
	"crypto/tls"
	"database/sql"
	"os"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq" //import postgres

	"github.com/go-pg/pg/v9"
)

//DB ...
type DB struct {
	*sql.DB
}

const (
	//DbUser ...
	DbUser = "postgres"
	//DbPassword ...
	DbPassword = "postgres"
	//DbName ...
	DbName = "qiyetalk_deveopment"
)

var db *gorp.DbMap

func getTLSConfig() *tls.Config {
	pgSSLMode := os.Getenv("PGSSLMODE")
	if pgSSLMode == "disable" {
		return nil
	}
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func pgOptions() *pg.Options {
	return &pg.Options{
		User:      DbUser,
		Password:  DbPassword,
		Database:  DbName,
		TLSConfig: nil,

		MaxRetries:      1,
		MinRetryBackoff: -1,

		DialTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		PoolSize:           10,
		MaxConnAge:         10 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        10 * time.Second,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}

var _db *pg.DB

// GetDB ...
func GetDB() *pg.DB {
	if _db == nil {
		_db = pg.Connect(pgOptions())
	}
	return _db
}
