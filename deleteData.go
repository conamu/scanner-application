package main

import (
	"encoding/csv"
	"github.com/spf13/viper"
	"fmt"
	"log"
	"os"
	"github.com/dgraph-io/badger/v2"
)

func deleteData(code string, newRecord []string, valid bool) {

	if !valid {
		return
	}

	if code == "" && newRecord == nil {
		return
	}

	file, err := os.OpenFile(viper.GetString("flatPath"), os.O_RDWR, 0755)
	defer file.Close()
	check(err)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)
	err = os.Remove(viper.GetString("flatPath"))
	check(err)

	nFile, err := os.OpenFile(viper.GetString("flatPath"), os.O_RDWR|os.O_CREATE, 0755)
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

func deleteBadger(code string, valid bool) {

	if !valid {
		return
	}

	err := bdb.Update(func(txn *badger.Txn) error {

		txn = bdb.NewTransaction(true)
		//deleting Name
		_, err := txn.Get([]byte(code + "Name"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This Item hasn't store in Database. You will be redirected to the main menu")
			sleep()
			completeMode()
			return nil
		} else if err != nil {
			log.Fatal(err)
		} else {
			//if there is a Name stored, it means, that other two values were also stored
			//even if with leer strings
			//so we can delete them, without checkinh, if there stored
			err = txn.Delete([]byte(code + "Name"))
			check(err)

			//deleting Category
			_, err = txn.Get([]byte(code + "Category"))
			check(err)
			err = txn.Delete([]byte(code + "Category"))
			check(err)

			//deleting Description
			_, err = txn.Get([]byte(code + "Description"))
			check(err)
			err = txn.Delete([]byte(code + "Description"))
			check(err)

			fmt.Println("Item is sucsessfuly deleted. You will be redirected to the main menu")
			txn.Commit()
			sleep()
			return nil
		}
		return nil
	})
	check(err)

}
