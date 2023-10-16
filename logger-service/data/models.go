package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs")
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})

	cur, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding logs")
		return nil, err
	}
	defer cur.Close(ctx)

	var logs []*LogEntry

	for cur.Next(ctx) {
		var item LogEntry
		err = cur.Decode(item)
		if err != nil {
			log.Println("error decoding log into slice")
			return nil, err
		}
		logs = append(logs, &item)
	}
	return logs, nil
}

func (l *LogEntry) GetById(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error decoding logId to docId")
		return nil, err
	}
	opts := options.FindOne()

	var item LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}, opts).Decode(item)
	if err != nil {
		log.Println("error decoding log while getting from DB")
		return nil, err
	}
	return &item, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("error dropping log")
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("error decoding logId to docId in update")
		return nil, err
	}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: l.Name},
		{Key: "data", Value: l.Data},
		{Key: "updatedAt", Value: time.Now()},
	}}})
	if err != nil {
		log.Println("error updating log entry")
		return nil, err
	}
	return result, nil
}
