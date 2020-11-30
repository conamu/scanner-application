package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/dgraph-io/badger/v2"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var notFound = errors.New("CODE NOT FOUND")

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

		return true, nil, notFound
	}

	if option != "5" && option != "6" {
		sleep()
	}

	return true, row, nil
}






func readKV(code string, valid bool) (bool, []string, error) {

	if !valid {
		return true, nil, nil
	}

	if code != "end" {
		exists := checkItem(code)
		if !exists {
			return true, nil, notFound
		}

		if viper.GetBool("apiEndpointMode") {
			result := make([]string, 4)
			result[0] = code
			result[1] = string(readName(code))
			result[2] = string(readCategory(code))
			result[3] = string(readDescription(code))

			return true, result, nil

		} else if !viper.GetBool("apiEndpointMode") {
			name := readName(code)
			category := readCategory(code)
			description := readDescription(code)
			itemDisplay(string(name), string(category), string(description))

			return true, nil, nil
		}
	}
	log.Println("Scanned end code, exiting!")
	sleep()
	return false, nil, nil

}

//this function is checking, if the Item already stored in the DB
func checkItem(code string) bool {
	err := bdb.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This Item hasn't store in Database.")
			return badger.ErrKeyNotFound
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

func readName(code string) []byte {
	//creating a copy of item, so we can use it later outside of transaction
	var nameCopy []byte
	err := bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(code + "Name"))
		check(err)

		err = item.Value(func(val []byte) error {
			//storing the value of item to the copy
			nameCopy = append([]byte{}, val...)
			return nil
		})
		check(err)
		return nil
	})
	check(err)
	return nameCopy

}

func readCategory(code string) []byte {
	//creating a copy of item, so we can use it later outside of transaction
	var catCopy []byte
	err := bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(code + "Category"))
		check(err)

		err = item.Value(func(val []byte) error {
			//storing the value of item to the copy
			catCopy = append([]byte{}, val...)
			return nil
		})
		check(err)
		return nil
	})
	check(err)
	return catCopy
}

func readDescription(code string) []byte {
	//creating a copy of item, so we can use it later outside of transaction
	var desCopy []byte
	err := bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(code + "Description"))
		check(err)
		err = item.Value(func(val []byte) error {
			//storing the value of item to the copy
			desCopy = append([]byte{}, val...)
			return nil
		})
		check(err)
		return nil
	})
	check(err)
	return desCopy
}
