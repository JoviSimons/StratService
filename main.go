package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/S-A-RB05/StratService/messaging"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	handleRequests()
}

//controllers

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAll")
	json.NewEncoder(w).Encode(getAllStrats())
}

func returnStrat(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	json.NewEncoder(w).Encode(readSingleStrat(idParam))
}

func useStrat(w http.ResponseWriter, r *http.Request) {
	// create a new Strategy object
	s := Strategy{}

	// retrieve the document with the specified _id and assign its values to the fields of the Strategy object
	var idParam string = mux.Vars(r)["id"]
	result := readSingleStrat(idParam)

	resultBytes, err := bson.Marshal(result)
	if err != nil {
		panic(err)
	}
	err = bson.Unmarshal(resultBytes, &s)
	if err != nil {
		panic(err)
	}

	fmt.Println("USING STRAT:")
	fmt.Println(s.Name)
	// Send script using rabbitmq
	//messaging.ProduceMessage(s.Mq, "strat_queue")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAll)
	myRouter.HandleFunc("/get/{id}", returnStrat)
	myRouter.HandleFunc("/use/{id}", useStrat)
	myRouter.HandleFunc("/create", storeStrat)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func storeStrat(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Storing Strat")
	// parse the request body into a Strategy struct
	var strat Strategy
	err := json.NewDecoder(body).Decode(&strat)
	fmt.Println(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// insert the strategy into the database
	insertStrat(strat, w)

	// Marshal the struct into a byte slice
	bStrat, err := json.Marshal(strat)
	if err != nil {
		panic(err)
	}

	messaging.ProduceMessage(bStrat, "q.syncStrat")

}

//service functions

func getAllStrats() (values []primitive.M) {
	return readAllStrats()
}

// other
// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
