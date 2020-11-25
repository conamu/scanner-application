package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

// Struct for JSON data structure
type CodeRes struct {
	Code string `json:code`
	Name string `json:name`
	Category string `json:category`
	Description string `json:description`
}

// This starts up the webserver on the configured port and accepts requests.
func requestHandler() {
	router := mux.NewRouter().StrictSlash(true)
	port := ":" + viper.GetString("apiEndpointPort")
	router.HandleFunc("/{code}", getCodeData)
	fmt.Printf("Barcode-Database Read-only API Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// This is the main fucntion that uses our built in functions to validate barcodes and, if valid, get Data from our Data backends.
// If not valid, it will return a 406 - Not Acceptable status code. If not found, it will return a 404 - Not Found status code.
func getCodeData(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: getCodeData")
	inputCode := mux.Vars(r)["code"]
	code, valid := validateBarcode(inputCode)

	if !valid {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("406 - Code not Valid."))
		return
	}

	if viper.GetBool("useFlatDB") {
		fmt.Println("FlatDB Read-Only mode")

		_, result, err := csvRead(code, "5", valid)
		if errors.Is(err, notFound){
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Code not Found."))
			return
		}

		res := CodeRes{
			Code:        strings.TrimSpace(result[0]),
			Name:        strings.TrimSpace(result[1]),
			Category:    strings.TrimSpace(result[2]),
			Description: strings.TrimSpace(result[3]),
		}

		json.NewEncoder(w).Encode(res)
		return

	} else if viper.GetBool("useKeyValueDB") {
		fmt.Println("Key-Value Database Read-only mode")
		check(db.View(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(code + "Name"))
			if err == badger.ErrKeyNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Code not Found"))
				return nil
			}

			nameVal, err := txn.Get([]byte(code+"Name"))
			categoryVal, err := txn.Get([]byte(code+"Category"))
			descriptionVal, err := txn.Get([]byte(code+"Description"))
			check(err)

			name, err := nameVal.ValueCopy(nil)
			category, err := categoryVal.ValueCopy(nil)
			description, err := descriptionVal.ValueCopy(nil)
			check(err)

			res := CodeRes{
				Code:        code,
				Name:        strings.TrimSpace(string(name)),
				Category:    strings.TrimSpace(string(category)),
				Description: strings.TrimSpace(string(description)),
			}

			json.NewEncoder(w).Encode(res)
			return nil

		}))
	}

}
