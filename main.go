package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/S-A-RB05/StratService/messaging"
	"github.com/gorilla/mux"
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

func getAllStratsOfUser(w http.ResponseWriter, r *http.Request) {
	var useridParam string = mux.Vars(r)["userid"]
	result, err := SearchStrats(useridParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAll)
	myRouter.HandleFunc("/get/{id}", returnStrat)
	myRouter.HandleFunc("/create", storeStrat)
	myRouter.HandleFunc("/getall/{userid}", getAllStratsOfUser)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func storeStrat(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Storing Strat")
	// parse the request body into a Strategy struct
	var strat Strategy
	err := json.NewDecoder(body).Decode(&strat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	var strat_id = insertStrat(strat)
	// insert the strategy into the database
	fmt.Fprint(w, strat_id)

	sendStratToTestManager(strat_id, strat)
}

//service functions

func getAllStrats() (values []primitive.M) {
	return readAllStrats()
}

func sendStratToTestManager(id string, strat Strategy) {
	var rStrat StrategyRequest

	rStrat.Id = id
	rStrat.Ex = strat.Ex
	rStrat.Name = strat.Name
	rStrat.Created = strat.Created

	// Marshal the struct into a byte slice
	bStrat, err := json.Marshal(rStrat)
	if err != nil {
		panic(err)
	}

	messaging.ProduceMessage(bStrat, "q.syncStrat")
}

// other
// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
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
