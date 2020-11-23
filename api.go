package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
)

type CodeRes struct {
	Code string `json:code`
	Name string `json:name`
	Category string `json:category`
	Description string `json:description`
}

func requestHandler() {
	router := mux.NewRouter().StrictSlash(true)
	port := ":" + viper.GetString("apiEndpointPort")
	router.HandleFunc("/{code}", getCodeData)
	fmt.Printf("Barcode-Database Read-only API Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

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
		file, err := os.OpenFile("data/database.csv", os.O_RDONLY|os.O_CREATE, 0755)
		check(err)
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()

		for _, record := range records {
			if record[0] == code {
				res := CodeRes{
					Code:        strings.TrimSpace(record[0]),
					Name:        strings.TrimSpace(record[1]),
					Category:    strings.TrimSpace(record[2]),
					Description: strings.TrimSpace(record[3]),
				}
				json.NewEncoder(w).Encode(res)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Code not Found."))
		return

	} else if viper.GetBool("useKeyValueDB") {

	}

}
