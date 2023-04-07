package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DbConfig struct {
	// https://gorm.io/docs/connecting_to_the_database.html
	DbDriver string
	DbName   string
	Host     string
	Port     int
	Username string
	Password string
	Timezone string
	// https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func buildConnStr(config DbConfig) (dsn string, err error) {
	// default timezone to lacal timezone
	if len(config.Timezone) == 0 {
		config.Timezone = time.Now().Location().String()
	}
	switch config.DbDriver {
	case "sqlite":
		// if empty db_name just use in memory database
		if len(config.DbName) == 0 {
			config.DbName = "file::memory:?cache=shared"
		}
		dsn = fmt.Sprintf("%s?cache=shared", config.DbName)
	case "mysql", "mariadb", "tidb":
		// if unix socket avaliavle use it
		_, err := os.Stat(config.Host)
		if errors.Is(err, os.ErrNotExist) {
			// tcp connection
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=true&loc=%s", config.Username, config.Password, config.Host, config.Port, config.DbName, config.Timezone)
		} else {
			// unix socket
			dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4,utf8&parseTime=true&loc=%s", config.Username, config.Password, config.Host, config.DbName, config.Timezone)
		}
	case "postgres", "pgx":
		dsn = fmt.Sprintf("config.Username=%s config.Passwordword=%s config.Host=%s config.Port=%d dbname=%s TimeZone=%s", config.Username, config.Password, config.Host, config.Port, config.DbName, config.Timezone)
	case "sqlserver":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", config.Username, config.Password, config.Host, config.Port, config.DbName)
	case "clickhouse":
		dsn = fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s", config.Host, config.Port, config.DbName, config.Username, config.Password)
	default:
		err = errors.New("not support DB_DRIVER")
		return
	}
	return
}

func ConnectDb(config DbConfig, opts ...gorm.Option) (dbConn *gorm.DB, err error) {
	dsn, err := buildConnStr(config)
	if err != nil {
		return
	}
	switch config.DbDriver {
	case "sqlite":
		dbConn, err = gorm.Open(sqlite.Open(dsn), opts...)
	case "mysql", "mariadb", "tidb":
		dbConn, err = gorm.Open(mysql.Open(dsn), opts...)
	case "postgres", "pgx":
		dbConn, err = gorm.Open(postgres.Open(dsn), opts...)
	case "sqlserver":
		dbConn, err = gorm.Open(sqlserver.Open(dsn), opts...)
	case "clickhouse":
		dbConn, err = gorm.Open(clickhouse.Open(dsn), opts...)
	default:
		err = errors.New("not support DB_DRIVER")
		return
	}
	if err != nil {
		return
	}
	err = applyDbPoolConfig(dbConn, config)
	return
}

func applyDbPoolConfig(dbConn *gorm.DB, config DbConfig) (err error) {
	// Get generic database object sql.DB to use its functions
	sqlDB, err := dbConn.DB()
	if err != nil {
		return
	}

	// set default connection pool
	// https://mariadb.com/docs/xpand/connect/programming-languages/nodejs/promise/connection-pools/#Options
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 10
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = config.MaxOpenConns / 2
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 30 * time.Minute
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	return
}
