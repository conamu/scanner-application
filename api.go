package main

import (
	"encoding/json"
	"errors"
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
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Endpoint Hit: getCodeData")
	inputCode := mux.Vars(r)["code"]
	if inputCode == "end" {
		if viper.GetBool("useKeyValueDB") {
			bdb.Close()
		}
		os.Exit(0)
	}
	code, valid := validateBarcode(inputCode)

	if !valid {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("406 - Code not Valid."))
		return
	}

	if viper.GetBool("useFlatDB") {
		fmt.Println("FlatDB Read-Only mode")

		_, record, err := csvRead(code, "5", valid)
		if errors.Is(err, notFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Code not Found."))
			return
		}

		res := CodeRes{
			Code:        strings.TrimSpace(record[0]),
			Name:        strings.TrimSpace(record[1]),
			Category:    strings.TrimSpace(record[2]),
			Description: strings.TrimSpace(record[3]),
		}
		json.NewEncoder(w).Encode(res)




	} else if viper.GetBool("useKeyValueDB") {
		fmt.Println("Key-Value Database Read-only mode")
		initDB()
		_, result, err := readKV(code, valid)

		if errors.Is(err, notFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Code not Found."))
			check(bdb.Close())
			return
		}

		check(bdb.Close())
		res := CodeRes{
			Code:        result[0],
			Name:        strings.TrimSpace(result[1]),
			Category:    strings.TrimSpace(result[2]),
			Description: strings.TrimSpace(result[3]),
		}

		json.NewEncoder(w).Encode(res)

	} else if viper.GetBool("useMysqlDB") {
		fmt.Println("MySql Database Read-only mode")
		initDB()
		_, result, err := readSql(code, valid)

		if errors.Is(err, notFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Code not Found."))
			check(mdb.Close())
			return
		}

		res := CodeRes{
			Code:        result[0],
			Name:        strings.TrimSpace(result[1]),
			Category:    strings.TrimSpace(result[2]),
			Description: strings.TrimSpace(result[3]),
		}
		check(mdb.Close())

		json.NewEncoder(w).Encode(res)
	}

}
