package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chat struct {
	ID        *uint64    `bson:"ID"`
	Messages  []*Message `bson:"messages"`
	CreatedAt *time.Time `bson:"createdAt"`
}

type Message struct {
	Data      *string    `bson:"data"`
	CreatedAt *time.Time `bson:"createAt"`
}

type ChatDAO interface {
	InsertOne(Message *Chat) (*mongo.InsertOneResult, error)
	FindByID(ID uint64) (*Chat, error)
}

func NewMongoDAO(uri, dbName, collectionName string) (*MongoDAO, error) {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Printf("Connected to MongoDB")
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDAO{
		client:     client,
		collection: collection,
	}, nil
}

type MongoDAO struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (dao MongoDAO) InsertOne(Message *Chat) (*mongo.InsertOneResult, error) {
	ctx := context.Background()
	result, err := dao.collection.InsertOne(ctx, Message)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

func (dao MongoDAO) FindByID(ID uint64) (*Chat, error) {
	ctx := context.Background()
	filter := bson.M{"ID": ID}
	cursor, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	var chat *Chat
	if cursor.Next(ctx) {
		if err := cursor.Decode(chat); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}
	}

	return chat, nil
}

func (dao MongoDAO) Close() {
	ctx := context.Background()
	if err := dao.client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	fmt.Println("Disconnected from MongoDB!")
}
