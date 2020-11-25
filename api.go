package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Article ...
type Article struct {
	Barcode     string `json:"barcode"`
	Name        string `json:"name"`
	Category    string `json:"category`
	Description string `json:"description"`
}

//slice to store data from DB
var Articles []Article

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	if viper.GetBool("useFlatDB") {
		myRouter.HandleFunc("/articles/{barcode}", returnSingleArticleFlat).Methods("GET")
	} else if viper.GetBool("useKeyValueDB") {
		myRouter.HandleFunc("/articles/{barcode}", returnSingleArticleKV).Methods("GET")
	}
	//you can also use it without log.Fatal, but I don’t know the difference
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnSingleArticleKV(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//taking the last part of URL, which is the barcode of the item, we are looking for
	vars := mux.Vars(r)
	key := vars["barcode"]

	//put information, that we get from DB, to a slice
	_, Articles = readKV(key, true)
	//checking, if item is stored in DB, if not, giving back an error
	if checkItem(key) == false {
		http.NotFound(w, r)
	} else {
		//here suppose to be also an information about HTTP Status, but I don’t know how to get it (yet)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Articles)

	}

}

func returnSingleArticleFlat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//taking the last part of URL, which is the barcode of the item, we are looking for
	vars := mux.Vars(r)
	key := vars["barcode"]
	_, article, _ := csvRead(key, "1", true)
	if len(article) == 0 {
		http.NotFound(w, r)
	} else {
		Articles = []Article{
			{Barcode: key, Name: strings.TrimSpace(article[0]), Category: strings.TrimSpace(article[1]), Description: strings.TrimSpace(article[2])}}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Articles)

	}
}
