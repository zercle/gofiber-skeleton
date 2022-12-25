package datasources

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	// Import postgres driver
	"gorm.io/driver/postgres"
)

// PostgreSQL for init connection
type PostgreSQLConfig struct {
	// Database connection
	connStr string

	// Database name
	DbName string

	// Optional.
	Username, Password string

	// Host of the PostgreSQL instance.
	//
	// If set, UnixSocket should be unset.
	Host string

	// Port of the PostgreSQL instance.
	//
	// If set, UnixSocket should be unset.
	Port string

	// UnixSocket is the filepath to a unix socket.
	//
	// If set, Host and Port should be unset.
	UnixSocket string

	// Set max idle connection at times
	MaxIdleConns int

	// Set max open connection at time
	MaxOpenConns int

	// Set connection life time
	ConnMaxLifetime time.Duration

	// Let's sql driver parse time
	ParseTime bool
}

// postgresDStoreString returns a connection string suitable for sql.Open.
func (c PostgreSQLConfig) postgresDStoreString() string {

	// Ensure postgres port
	if port, err := strconv.Atoi(c.Port); err != nil || port > 65535 || port < 1 {
		c.Port = "9920"
	}

	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	if c.UnixSocket != "" {
		if _, err := os.Stat(c.UnixSocket); err == nil {
			return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, c.DbName)
		}
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok", c.Host, c.Username, c.Password, c.DbName, c.Port, "disable")
}

// New PostgreSQL creates a new database connection backed by a given postgres server.
func (config PostgreSQLConfig) InitPostgreSqlConn(dbName string) (dbConn *gorm.DB, err error) {
	if len(dbName) == 0 {
		return nil, fmt.Errorf("need: dbName")
	}

	config.DbName = dbName

	config.connStr = config.postgresDStoreString()
	// +07:00
	config.connStr = config.connStr + "?loc=Asia%2FBangkok&time_zone=%27%2B07%3A00%27"
	// Asia/Bangkok
	// conn.ConnStr = conn.ConnStr + "?loc=Asia%2FBangkok&time_zone=%27Asia%2FBangkok%27"
	config.connStr = config.connStr + "&charset=utf8mb4,utf8"
	if config.ParseTime {
		config.connStr = config.connStr + "&parseTime=true"
	}

	// Use system default database if empty
	if len(config.connStr) == 0 {
		return nil, fmt.Errorf("PostgreSQL: connStr needed")
	}
	// Open connection to database
	dbConn, err = gorm.Open(postgres.Open(config.connStr), &gorm.Config{
		PrepareStmt: true,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Printf("NewPostgreSQLDB: \n%+v", err)
		return nil, fmt.Errorf("PostgreSQL: could not get a connection: %v", err)
	}

	err = config.ApplyConnOption(dbConn)
	if err != nil {
		log.Printf("NewPostgreSQLDB: \n%+v", err)
		return nil, fmt.Errorf("PostgreSQL: could not config connection: %v", err)
	}

	return
}

// ApplyConnOption to current db connection
func (c *PostgreSQLConfig) ApplyConnOption(dbConn *gorm.DB) (err error) {
	dbObj, err := dbConn.DB()

	// Set max open connection at time
	if c.MaxOpenConns != 0 {
		dbObj.SetMaxOpenConns(c.MaxOpenConns)
	} else {
		// Default value follow PostgreSQL.js pool
		dbObj.SetMaxOpenConns(10)
	}

	// Set max idle connection at time
	if c.MaxIdleConns != 0 {
		dbObj.SetMaxIdleConns(c.MaxIdleConns)
	} else {
		// Default value follow PostgreSQL.js pool
		dbObj.SetMaxIdleConns(5)
	}

	// Time out for long
	if c.ConnMaxLifetime != 0 {
		dbObj.SetConnMaxLifetime(c.ConnMaxLifetime)
	} else {
		// Default value follow PostgreSQL.js pool
		dbObj.SetConnMaxLifetime(5 * time.Minute)
	}

	return
}
