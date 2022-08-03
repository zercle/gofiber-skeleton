package datasources

import (
	"fmt"
	"log"

	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"
	// Import mongo driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig for init connection
type MongoDBConfig struct {
	// Database connection
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
	Port string
}

// New MongoDB creates a new database connection backed by a given mongodb server.
func (config MongoDBConfig) NewMongoDB(dbname string) (dbConn *mongo.Database, err error) {
	// Use system default database if empty
	if dbname == "" {
		dbname = viper.GetString("db.main.db_name")
	}

	// Create new client
	config.ConnStr = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Username, config.Password, config.Host, config.Port)
	client, err := mgm.NewClient(options.Client().ApplyURI(config.ConnStr))
	if err != nil {
		log.Printf("NewMongoDB: \n%+v", err)
		return nil, fmt.Errorf("MongoDB: could not get a connection: %v", err)
	}

	// Get the model's db
	dbConn = client.Database(dbname)

	return
}
