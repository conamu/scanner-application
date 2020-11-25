package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/badger/v2"
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

//not sure we need it, but this slice is to simulate a database
var Articles []Article

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	if viper.GetBool("useFlatDB") {
		myRouter.HandleFunc("/articles/{barcode}", returnSingleArticleFlat).Methods("GET")
	} else if viper.GetBool("useKeyValueDB") {
		//add articles route and map it to responsible function
		//myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
		myRouter.HandleFunc("/articles/{barcode}", returnSingleArticleKV).Methods("GET")
		//you can also use it without log.Fatal, but I don’t know the difference

	}
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

//if I will have enough time, I would also try this one to work
/* func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	//make the result look nicer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Articles)
} */

func returnSingleArticleKV(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//taking the last part of URL, which is the barcode of the item, we are looking for
	vars := mux.Vars(r)
	key := vars["barcode"]

	//creating three variables, where we gonna store relevant data
	var name string
	var category string
	var description string

	//Open KV DB and loop over it to find a match items
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		//not sure, if we need it, read about in BD docs
		//opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				if string(k) == key+"Name" {
					name = string(v)
				} else if string(k) == key+"Category" {
					category = string(v)
				} else if string(k) == key+"Description" {
					description = string(v)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	//put information, that we get from DB, to a slice of type Article
	Articles = []Article{
		{Barcode: key, Name: strings.TrimSpace(name), Category: strings.TrimSpace(category), Description: strings.TrimSpace(description)},
	}
	//if we created a slice with not stored Barcode, we should give back an error
	if len(name) == 0 && len(category) == 0 && len(description) == 0 {
		http.NotFound(w, r)
	} else {
		//here suppose to be also an information about HTTP Status, but I don’t know how to get it (yet)
		json.NewEncoder(w).Encode(Articles)
		fmt.Println(Articles)
	}
	check(err)
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

		json.NewEncoder(w).Encode(Articles)
		fmt.Println(Articles)
	}
}
