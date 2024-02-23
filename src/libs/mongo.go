package libs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Annotation struct {
	ID          primitive.ObjectID `bson:"_id"`
	Manifest_id string             `bson:"manifest_id"`
	Chars       string             `bson:"chars"`
	UserId      string             `bson:"userId"`
	CreatedAt   string             `bson:"createdAt"`
	Motivation  string             `bson:"motivation,omitempty"`
	On          string             `bson:"on"`
}

func DBConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	uri := os.Getenv("DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB Connected")
	return client
}

func GetAnnotations(target string, motivation string, user string, date string, client *mongo.Client) []Annotation {
	filter := bson.D{{Key: "chars", Value: target}, {Key: "motivation", Value: motivation}, {Key: "user", Value: user}, {Key: "date", Value: date}}
	opt := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})

	annotations := client.Database("test").Collection("annotations")
	cursor, err := annotations.Find(context.TODO(), filter, opt)
	if err != nil {
		panic(err)
	}

	var results []Annotation
	for cursor.Next(context.TODO()) {
		var result bson.D
		var annotation Annotation
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		doc, err := bson.Marshal(result)
		if err != nil {
			panic(err)
		}
		if err := bson.Unmarshal(doc, &annotation); err != nil {
			panic(err)
		}
		results = append(results, annotation)
	}

	return results
}
