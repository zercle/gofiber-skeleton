package datasources

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
	// Import mysql driver
	"gorm.io/driver/mysql"
)

// DateTimeFormat Long date time mysql format
const DateTimeFormat = "2006-02-01 15:04:05"

// DateFormat Date mysql format
const DateFormat = "2006-02-01"

// TimeFormat Time mysql format
const TimeFormat = "15:04:05"

// MariadbConfig for init connection
type MariadbConfig struct {
	// Database connection
	ConnObj *gorm.DB
	ConnStr string

	// Database name
	DBName string

	// Optional.
	Username, Password string

	// Host of the mariadb instance.
	//
	// If set, UnixSocket should be unset.
	Host string

	// Port of the mariadb instance.
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

// mariadbDStoreString returns a connection string suitable for sql.Open.
func (c MariadbConfig) mariadbDStoreString() string {

	// Ensure mariadb port
	if port, err := strconv.Atoi(c.Port); err != nil || port > 65535 || port < 1 {
		c.Port = "3306"
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
		return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, c.DBName)
	}
	return fmt.Sprintf("%stcp([%s]:%s)/%s", cred, c.Host, c.Port, c.DBName)
}

// NewMariadbDB creates a new database connection backed by a given mariadb server.
func (c MariadbConfig) NewMariadbDB(dbname string) (conn *MariadbConfig, err error) {
	// Use system default database if empty
	if len(dbname) == 0 {
		dbname = os.Getenv("DB_NAME")
	}

	// Init conn
	conn = &c

	conn.DBName = dbname

	conn.ConnStr = conn.mariadbDStoreString()
	// +07:00
	conn.ConnStr = conn.ConnStr + "?loc=Asia%2FBangkok&time_zone=%27%2B07%3A00%27"
	// Asia/Bangkok
	// conn.ConnStr = conn.ConnStr + "?loc=Asia%2FBangkok&time_zone=%27Asia%2FBangkok%27"
	conn.ConnStr = conn.ConnStr + "&charset=utf8mb4,utf8"
	if conn.ParseTime {
		conn.ConnStr = conn.ConnStr + "&parseTime=true"
	}

	// Use system default database if empty
	if len(conn.ConnStr) == 0 {
		return nil, fmt.Errorf("MariaDB: connStr needed")
	}
	// Open connection to database
	conn.ConnObj, err = gorm.Open(mysql.Open(conn.ConnStr), &gorm.Config{
		PrepareStmt: true,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Printf("NewMariadbDB: \n%+v", err)
		return nil, fmt.Errorf("MariaDB: could not get a connection: %v", err)
	}

	err = conn.ApplyConnOption()
	if err != nil {
		log.Printf("NewMariadbDB: \n%+v", err)
		return nil, fmt.Errorf("MariaDB: could not config connection: %v", err)
	}

	return
}

// ApplyConnOption to current db connection
func (c *MariadbConfig) ApplyConnOption() (err error) {
	dbObj, err := c.ConnObj.DB()

	// Set max open connection at time
	if c.MaxOpenConns != 0 {
		dbObj.SetMaxOpenConns(c.MaxOpenConns)
	} else {
		// Default value follow mariadb.js pool
		dbObj.SetMaxOpenConns(10)
	}

	// Set max idle connection at time
	if c.MaxIdleConns != 0 {
		dbObj.SetMaxIdleConns(c.MaxIdleConns)
	} else {
		// Default value follow mariadb.js pool
		dbObj.SetMaxIdleConns(5)
	}

	// Time out for long
	if c.ConnMaxLifetime != 0 {
		dbObj.SetConnMaxLifetime(c.ConnMaxLifetime)
	} else {
		// Default value follow mariadb.js pool
		dbObj.SetConnMaxLifetime(5 * time.Minute)
	}

	return
}

// GetCurrentConn get db connection
func (c *MariadbConfig) GetCurrentConn() (conn *gorm.DB, err error) {
	dbObj, err := c.ConnObj.DB()
	checkPing := dbObj.Ping()
	if checkPing != nil {
		log.Printf("GetCurrentConn: \n%+v", checkPing)
		c.ConnObj, err = gorm.Open(mysql.Open(c.ConnStr), &gorm.Config{PrepareStmt: true})
		if err != nil {
			log.Printf("GetCurrentConn: \n%+v", err)
			return nil, fmt.Errorf("MariaDB: could not get a connection: %v", err)
		}
		err = c.ApplyConnOption()
		if err != nil {
			log.Printf("GetCurrentConn: \n%+v", err)
			return nil, fmt.Errorf("MariaDB: could not config connection: %v", err)
		}
	}
	return c.ConnObj, err
}

// Close close current db connection.  If database connection is not an io.Closer, returns an error.
func (c *MariadbConfig) Close() (err error) {
	dbObj, err := c.ConnObj.DB()
	if err != nil {
		return err
	}
	err = dbObj.Close()
	return
}
