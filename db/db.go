package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"polo.gamesmania.io/custodia/gp"
)

const DefaultUserID = "3b96b8ae-991e-4ccc-b703-027612295616"

// Ctx is the mongo context
var Ctx context.Context

// LocalDB is the mongo db
var DBClient *mongo.Client
var LocalDB *mongo.Database

const ListLimit = 100

// Connect connects to mongoDB
func Connect() error {

	// Set client options
	Ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(gp.PConfig.MongoDBURI)

	// Connect to MongoDB
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		return err
	}

	// We make sure we have been connected
	err = client.Ping(Ctx, readpref.Primary())
	if err != nil {
		return err
	}

	DBClient = client
	LocalDB = DBClient.Database(gp.PConfig.DBName)

	return nil
}

// ReconnectCheck reconnects to DB
func ReconnectCheck() {

	// We make sure we are still connected
	err := LocalDB.Client().Ping(Ctx, readpref.Primary())
	if err == nil {
		return
	}

	// We reconnect
	Connect()
}

// Init inits default features we need for the database to work
func Init() error {

	err := CreateIndexes()
	if err != nil {
		return err
	}

	return nil
}

func CreateIndexes() error {

	// TODO: Create indexes here
	return nil
}
