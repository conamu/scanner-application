package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	//myRouter.HandleFunc("/", homePage)
	//add articles route and map it to responsible function
	//myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles/{barcode}", returnSingleArticle).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

/* func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: Homepage")

} */

/* func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	//make the result look nicer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Articles)
}
*/
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//vars := mux.Vars(r)
	//key := vars["barcode"]

	json.NewEncoder(w).Encode(Articles)

}

/*
 func main() {
	fmt.Println("Rest API v 2.0 - Mux Routers")
	Articles = []Article{
		Article{Barcode: "1234", Name: "Cola", Category: "Drinks", Description: "Sweet soft drink"},
		Article{Barcode: "1235", Name: "Pommes", Category: "Snacks", Description: "Makes you fat. but happy"},
	}
	handleRequests()
}

*/
