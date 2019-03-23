package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookType struct {
	Name     string
	Genre    string
	ID       string `bson:"id" json:"id` //The automatically generated id's are not currently supported
	AuthorID string `bson:"authorId" json:"authorId`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, error := mongo.NewClient(options.Client().ApplyURI("ur_Database_uri"))
	error = client.Connect(ctx)

	//Checking the connection
	error = client.Ping(context.TODO(), nil)
	fmt.Println("Database connected")

	//Specify your respective collection
	BooksCollection := client.Database("test").Collection("books")

	/**
	 * Create - Adding a new book
	 * res -> the insert command returns the inserted id of the oject
	 */
	res, err := BooksCollection.InsertOne(ctx, bson.M{"name": "The Go Language", "genre": "Coding", "authorId": "4"})
	if err != nil {
		log.Fatal(err)
	}

	/**
	 * Read -
	 * Getting all the books back
	 */
	cur, error := BooksCollection.Find(ctx, bson.D{{}})

	var allbooks []*bookType

	for cur.Next(context.TODO()) {
		var bookHolder bookType
		err := cur.Decode(&bookHolder)
		if err != nil {
			log.Fatal(err)
		}
		allbooks = append(allbooks, &bookHolder)
	}
	defer cur.Close(context.TODO())
	for _, element := range allbooks {
		book := *element
		fmt.Println(book)
	}

	if error != nil {
		log.Fatal(error)
	}

	/**
	 * Update
	 * Collection has functions like UpdateOne and UpdateMany
	 * Returns the Matched and Modified Count
	 */
	filter := bson.D{{"name", "Book 1"}}
	newName := bson.D{
		{"$set", bson.D{
			{"name", "Updated Name of Book 1"},
		}},
	}

	res, err := BooksCollection.UpdateOne(ctx, filter, newName)
	if err != nil {
		log.Fatal(err)
	}
	updatedObject := *res
	fmt.Printf("The matched count is : %d, the modified count is : %d", updatedObject.MatchedCount, updatedObject.ModifiedCount)

	/**
	 * Delete
	 */
	filter = bson.D{{"name", "Book 2"}}
	deleteResult, error := BooksCollection.DeleteOne(ctx, filter)
	fmt.Println(deleteResult)
}
