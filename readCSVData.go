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

func readBadger(code string) {

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(code))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This entry haven't store yet. You will be redirected to the main menu")
			time.Sleep(time.Second * 4)
			main()
		} else if err != nil {
			log.Fatal(err)
		}
		err = item.Value(func(val []byte) error {
			fmt.Printf("The answer is: %s\n", val)
			time.Sleep(time.Second * 4)
			return nil
		})
		return nil
	})

	check(err)

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
