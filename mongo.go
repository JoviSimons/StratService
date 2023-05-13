package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// global variable mongodb connection client
var client mongo.Client = newClient()

// ----Create----
func insertStrat(strat Strategy) (writeValue string) {
	stratCollection := client.Database("testing").Collection("strategies")
	strat.Created = time.Now()
	result, err := stratCollection.InsertOne(context.TODO(), strat)
	if err != nil {
		panic(err)
	}

	// return the ID of the newly inserted script
	writeValue = result.InsertedID.(primitive.ObjectID).Hex()
	return writeValue
}

//----Read----

func readAllStrats() (values []primitive.M) {
	stratCollection := client.Database("testing").Collection("strategies")
	// retrieve all the documents (empty filter)
	cursor, err := stratCollection.Find(context.TODO(), bson.D{})
	// check for errors in the finding
	if err != nil {
		panic(err)
	}

	// convert the cursor result to bson
	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// display the documents retrieved
	fmt.Println("displaying all results from the search query")
	for _, result := range results {
		fmt.Println(result)
	}

	values = results
	return
}

func readSingleStrat(id string) (value primitive.M) {
	stratCollection := client.Database("testing").Collection("strategies")
	// convert the hexadecimal string to an ObjectID type
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	// retrieve the document with the specified _id
	var result bson.M
	err = stratCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&result)
	if err != nil {
		panic(err)
	}

	// display the retrieved document
	fmt.Println("displaying the result from the search query")
	fmt.Println(result)
	value = result

	return value
}

//----Update----

//----Delete----

// other
func newClient() (value mongo.Client) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://stockbrood:admin@stockbrood.sifn3lq.mongodb.net/test"))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	value = *client

	return
}
