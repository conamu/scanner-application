package main

import (
	"encoding/csv"
	"fmt"
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

func readKV(code string) bool {
	checkItem(code)
	if code != "end" {
		name := readName(code)
		category := readCategory(code)
		description := readDescription(code)
		itemDisplay(string(name), string(category), string(description))

		return true
	}
	log.Println("Scanned end code, exiting!")
	return false
}

//this function is checking, if the Item already stored in the DB

/* func checkItem(code string, option int) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {
				fmt.Println("This Item haven't store in Database. You will be redirected to the main menu")
				main()
			} else if err != nil {
				log.Fatal(err)
				return err
			}
			return nil

		}
		check(err)
		return nil
	})
	if err != nil {
		return false
	}
	return true
} */

/* func checkItem1(code string) bool {
	var option bool
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {

			option = false
			return nil

		} else if err != nil {
			log.Fatal(err)
			return err
		}

	})
	if err != nil || option == false {
		return false
	}
	return true
} */

func checkItem(code string) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This Item haven't store in Database. You will be redirected to the main menu")
			time.Sleep(time.Second * 4)
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
