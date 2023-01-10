package datasources

import (
	"fmt"
	"log"

	"github.com/kamva/mgm/v3"

	// Import mongo driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig for init connection
type MongoDBConfig struct {
	// Database connection
	connStr string

	// Database name
	DbName string

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
func (config *MongoDBConfig) InitMongoDbConn(dbName string) (dbConn *mongo.Database, err error) {
	if len(dbName) == 0 {
		return nil, fmt.Errorf("need: dbName")
	}

	// config.DbName = dbName

	// Create new client
	config.connStr = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Username, config.Password, config.Host, config.Port)
	client, err := mgm.NewClient(options.Client().ApplyURI(config.connStr))
	if err != nil {
		log.Printf("NewMongoDB: \n%+v", err)
		return nil, fmt.Errorf("MongoDB: could not get a connection: %v", err)
	}

	// Get the model's db
	dbConn = client.Database(dbName)

	return
}
