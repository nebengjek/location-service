package mysql

import (
	"fmt"
	"location-service/bin/config"
	"location-service/bin/pkg/log"
	"location-service/bin/pkg/utils"
	"runtime"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DBConn variable to declare Database Connection
var DBConn *DatabaseConnection

type (
	// DBInterface to provide general Func
	DBInterface interface {
		Connect(string) *DatabaseConnection
		GetDB(string) (*sqlx.DB, error)
	}

	// Database to provide Database Config
	Database struct {
		Name DatabaseConfig
	}

	// DatabaseConfig currently have master only
	DatabaseConfig struct {
		Master string
	}

	// DatabaseConnection provide struct sqlx connection
	DatabaseConnection struct {
		Connection *sqlx.DB
	}
)

var (
	accessOnce sync.Once
	access     DBInterface
)

func InitConnection() DBInterface {
	if access != nil {
		return access
	}

	accessOnce.Do(func() {
		dsn, dbname := config.GetConfig().DnsMariaDB()
		dbClient := NewDatabase(dsn)
		dbClient.Connect(dbname)
		access = dbClient
	})

	return access
}

func StatusConnection() DBInterface {
	dsn, dbname := config.GetConfig().DnsMariaDB()
	dbClient := NewDatabase(dsn)
	dbClient.GetDB(dbname)

	return dbClient
}

// NewDatabase create Database Struct from Config
func NewDatabase(config interface{}) *Database {
	cfg := config.(string)

	dc := DatabaseConfig{
		Master: cfg,
	}

	return &Database{
		Name: dc,
	}
}

// Connect provide sqlx connection
func (db *Database) Connect(dbName string) *DatabaseConnection {

	databaseConn := DatabaseConnection{}

	master := db.Name.Master
	if master != "" {
		db, err := sqlx.Connect("mysql", master)
		if err != nil {
			log.GetLogger().Error("mysql", "Can not connect MySQL", "connect", utils.ConvertString(err))
		}

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(time.Minute * 10)

		databaseConn.Connection = db
	}

	DBConn = &DatabaseConnection{Connection: databaseConn.Connection}
	return DBConn
}

// GetDB provide status from DB
func (db *Database) GetDB(dbName string) (*sqlx.DB, error) {
	var newDB *sqlx.DB

	newDB = DBConn.Connection

	if newDB.Stats().OpenConnections > 40 {
		fpcs := make([]uintptr, 1)
		n := runtime.Callers(2, fpcs)
		if n != 0 {
			fun := runtime.FuncForPC(fpcs[0] - 1)
			if fun != nil {
				log.GetLogger().Error("mysql", fmt.Sprintf("Db Conn more than 40, Caller from Func : %s", fun.Name()), "GetDB", string(fun.Name()))
			}
		}
		log.GetLogger().Info("mysql", fmt.Sprintf("DB Conn more than 40, currently : %s", utils.ConvertString(newDB.Stats())), "GetDB", utils.ConvertString(newDB.Stats()))
	}
	return newDB, nil
}
