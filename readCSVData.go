package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/viper"
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

func csvRead(code string, option string, validity bool) (bool, []string, error) {

	if !validity {
		return false, nil, nil
	}
	var row []string
	file, err := os.OpenFile(viper.GetString("flatPath"), os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	check(err)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)
	notCount := 1

	// If the code matches an entry in the Database, show the data. Else return an error.
	for _, record := range records {
		if stringInSlice(code, record) {
			itemDisplay(record[1], record[2], record[3])
			row = record
		} else if code == "end" {
			log.Println("Scanned end code, exiting!")
			return false, row, nil
		} else if !stringInSlice(code, record) {
			notCount++
		}
	}

	if notCount > len(records) {
		fmt.Println("This code is not stored in the system.")

		return true, nil, errors.New("CODE NOT FOUND")
	}

	if option != "5" && option != "6" {
		time.Sleep(time.Second * 4)
	}

	return true, row, nil
}


