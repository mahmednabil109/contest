package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/linuxboot/contest/cmds/admin_server/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultDB                = "admin-server-db"
	DefaultCollection        = "logs"
	DefaultConnectionTimeout = 10 * time.Second
)

type MongoStorage struct {
	dbClient   *mongo.Client
	collection *mongo.Collection
}

func NewMongoStorage(uri string) (storage.Storage, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectionTimeout)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// check that the server is alive
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("Err while pinging mongo server: %w", err)
	}

	collection := client.Database(DefaultDB).Collection(DefaultCollection)

	return &MongoStorage{
		dbClient:   client,
		collection: collection,
	}, nil
}

func (s *MongoStorage) StoreLog(log storage.Log) error {
	logEntry := Log{
		ID:  primitive.NewObjectID(),
		Log: log,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectionTimeout)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, logEntry)

	return err
}

func (s *MongoStorage) GetLogs(opt storage.SearchOpt) ([]storage.Log, error) {
	// TODO
	return nil, nil
}

func (s *MongoStorage) Close() error {
	// TODO
	return nil
}

type Log struct {
	ID primitive.ObjectID `json:"ID"`
	storage.Log
}
