package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)


func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func csvRead(code string, option string) (bool, []string) {
	var row []string
	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	check(err)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)
	notCount := 1

	// If the code matches an entry in the Databease, show the data. Else return an error.
	for index, record := range records {
		if stringInSlice(code, record) {
			fmt.Println("====================================================\n",
				"Nr.: ", index,
				" == ", record[0],
				" == ", record[1],
				" == ", record[2],
				"\nDescription: ", record[3],

				"\n====================================================")
			row = record
		} else if code == "end" {
			log.Println("Scanned end code, exiting!")
			return false, row
		} else if !stringInSlice(code, record) {
			notCount++
		}
	}

	if notCount > len(records) {
		fmt.Println("This code is not stored in the system.")
	}

	if option != "5" && option != "6" {
		time.Sleep(time.Second * 4)
	}

return true, row
}


