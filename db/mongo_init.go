package db

import (
	"context"
	"os"
	"time"

	"ws-chat/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dBName         string
	collectionName string
	MongoClient    *mongo.Client
)

func InitMongo() error {
	// load mongo env parameters
	mongoURI := os.Getenv("MONGO_URI")
	dBName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	logger.Info("MONGO_URI: %s, DB_NAME: %s, COLLECTION_NAME: %s", mongoURI, dBName, collectionName)

	if mongoURI == "" || dBName == "" || collectionName == "" {
		logger.Error("[❌] MONGO_ENV parameters is not set! Please check environment variables.")
	}

	// connect to MongoDB
	mongoClient, err := createMongoConnection(mongoURI)
	if err != nil {
		logger.Error("[❌] MongoDB connect error: %v", err)
	}

	logger.Info("[✅] Create a Mongo connection to address: %s, DBname: %s, Collection: %s", mongoURI, dBName, collectionName)

	// check connection
	if err := CheckMongoConnection(mongoClient); err != nil {
		logger.Error("[❌] MongoDB connection failed: %v", err)
	}

	// init mongo collection
	if err := InitMongoCollection(mongoClient, dBName, collectionName); err != nil {
		logger.Error("[❌] MongoDB collection init failed: %v", err)
	}
	MongoClient = mongoClient
	return nil
}

// create a Mongo Connection
func createMongoConnection(mongoURI string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return mongoClient, nil
}

// Check Mongo Connection
func CheckMongoConnection(mongoClient *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mongoClient.Ping(ctx, nil)
}

// Init Mongo Collection, checks and creates collection if not exists, then rechecks connection
func InitMongoCollection(mongoClient *mongo.Client, dbName string, collectionName string) error {
	db := mongoClient.Database(dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	found := false
	for _, c := range collections {
		if c == collectionName {
			found = true
			break
		}
	}
	if !found {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			return err
		}
		logger.Info("[✅] Mongo messages collection created!")
	} else {
		logger.Info("[ℹ️] Mongo messages collection already exists!")
	}
	// recheck connection
	return CheckMongoConnection(mongoClient)
}
