package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"

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
		sleep()
	}

	return true, row, nil
}

func readKV(code string) bool {

	if code != "end" {
		checkItem(code)

		name := readName(code)
		category := readCategory(code)
		description := readDescription(code)
		itemDisplay(string(name), string(category), string(description))

		return true
	}
	log.Println("Scanned end code, exiting!")
	sleep()
	return false

}

//this function is checking, if the Item already stored in the DB
func checkItem(code string) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This Item hasn't store in Database. You will be redirected to the main menu")
      sleep()
			main()
			return nil
		} else if err != nil {
			log.Fatal(err)
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
	err := db.View(func(txn *badger.Txn) error {
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
	err := db.View(func(txn *badger.Txn) error {
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
	err := db.View(func(txn *badger.Txn) error {
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

//I NEEDED THIS FUNCTION TO CHECK, IF MY READ FUNCTION WORKS
//LET IT HERE JUST IN CASE
/* func addTest() error {

	fmt.Println("\nRunning SET")
	return db.Update(
		func(txn *badger.Txn) error {
			if err := txn.Set([]byte("lalala"), []byte("zhuzhuzhu")); err != nil {
				return err
			}
			fmt.Println("Set lalala to zhuzhuzhu")
			return nil
		})
}
*/
