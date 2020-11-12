package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v2"
)

func deleteData(code string, newRecord []string) {
	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR, 0755)
	defer file.Close()
	check(err)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)
	err = os.Remove("data/testDatabase.csv")
	check(err)

	nFile, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	defer nFile.Close()
	check(err)
	writer := csv.NewWriter(nFile)

	if code == "" {
		var indexToEdit int
		for index, record := range records {
			if record[0] == newRecord[0] {
				indexToEdit = index
			}
		}

		for i := 1; i < len(records[indexToEdit]); i++ {
			records[indexToEdit][i] = newRecord[i]
		}

		writer.WriteAll(records)

	} else {
		for _, record := range records {
			if record[0] != code {
				writer.Write(record)
			}
		}
	}
	writer.Flush()
}

/* err = db.Update(func(txn *badger.Txn) error {
txn := db.NewTransaction(true) // Read-write txn
defer txn.Discard()
err := txn.Delete([]byte(code))

check(err)
}
return nil */

func deleteBadger(code string) {
	err := db.Update(func(txn *badger.Txn) error {
		txn = db.NewTransaction(true)
		_, err := txn.Get([]byte(code + "Name"))
		check(err)
		err = txn.Delete([]byte(code + "Name"))
		check(err)
		fmt.Println("Item is sucsessfuly deleted. You will be redirected to the main menu")
		txn.Commit()
		time.Sleep(time.Second * 4)
		return nil

	})
	check(err)

}
