package client

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

type MongoDBService struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// Define the MongoDB collection and client as package-level variables
//var coll *mongo.Collection
//var client *mongo.Client

// NewMongoDBService creates a new MongoDBService instance.
func NewMongoDBService() (*MongoDBService, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGODB_URI not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully connect to DB\n")
	database := client.Database("ptt_data")
	//coll := database.Collection("gossip")
	return &MongoDBService{client: client, database: database}, nil
}

// Close closes the MongoDB client.
func (s *MongoDBService) Close() {
	if s.client != nil {
		_ = s.client.Disconnect(context.Background())
	}
}

// UpdateTopic updates the topic data in the database.
func (s *MongoDBService) UpdateTopic(boardName string, topics []Topic) (insertCount int, updateCount int) {
	collection := s.database.Collection(boardName)
	for _, topic := range topics {
		if topic.Title != "" && topic.URL != "" {
			filter := bson.M{"$and": []bson.M{
				{"title": topic.Title},
				{"url": topic.URL},
			}}
			var existingTopic Topic
			err := collection.FindOne(context.Background(), filter).Decode(&existingTopic)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					// Topic does not exist, insert new data
					_, err := collection.InsertOne(context.Background(), topic)
					if err != nil {
						log.Printf("Fail to insert data into database. err:%v\n", err)
					}
					insertCount++
					//fmt.Printf("New document inserted with ID: %s\n", insertResult.InsertedID)
				} else {
					// handle query error
					log.Fatalf("Fail to query data in database. err:%v\n", err)
				}
			} else {
				if existingTopic.Comments < topic.Comments {
					// Topic exists, update data
					update := bson.M{"$set": bson.M{
						"author":     topic.Author,
						"date":       topic.Date,
						"comments":   topic.Comments,
						"updatetime": topic.UpdateTime,
					}}
					_, err := collection.UpdateOne(context.Background(), filter, update)
					if err != nil {
						log.Fatalf("Fail to update data into database. err:%v\n", err)
					}
					//fmt.Printf("Document updated, id: %s, title: %s, existing comment: %d, new comments: %d\n", existingTopic.ID, existingTopic.Title, existingTopic.Comments, topic.Comments)
					updateCount++
				}

			}
		}
	}
	return insertCount, updateCount
}

func (s *MongoDBService) GetBoardDataFromDB(board string, counter int64) ([]Topic, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Set collection by board name
	s.collection = s.database.Collection(board)
	// Define the filter to query the last 'counter' number of data for the given board
	options := options.Find().SetSort(bson.M{"date": -1}).SetLimit(counter)

	// Perform the query
	cursor, err := s.collection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode documents
	var topics []Topic
	for cursor.Next(ctx) {
		var topic Topic
		if err := cursor.Decode(&topic); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	//fmt.Printf("Queried %d topics from MongoDB\n", len(topics))
	return topics, nil
}
