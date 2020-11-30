package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"
)

func writeData(data []string, barcode string, valid bool) bool {

	//open a file with flags: to append (O_Append) and to write(O_WRONLY)
	//FileMode (permission) - to append only
	file, err := os.OpenFile(viper.GetString("flatPath"), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)

	//creating a two-dimensional slice, where we save input-slices
	var products [][]string

	if len(data) == 0 {

		if !valid {
			return true
		}

		for _, record := range records {
			if record[0] == barcode {
				log.Println("The Barcode", barcode, " is already added in the System.")
				return true
			}
		}

		if barcode == "end" {
			return false
		}

		name, category, description := getParams()

		//creating a slice "product", which holds input as it's elements
		product := []string{barcode, name, category, description}

		//appending a slice "product" to two-dimensional slice "products"
		products = append(products, product)

	} else {

		if barcode == "end" {
			return false
		}

		fmt.Print("Whats the Product name? (max. 150 characters): ")
		scanner.Scan()
		name := scanner.Text()
		name = charLimiter(name, 150)

		for _, value := range records {

			if value[0] == data[0] {

				fmt.Print("Whats the Product description? (max. 500 characters): ")
				scanner.Scan()
				description := scanner.Text()
				description = charLimiter(description, 500)
				fmt.Println("====================================================")

				products = append(products, data)
				writer := csv.NewWriter(file)
				writer.WriteAll(products)
				if err := writer.Error(); err != nil {
					log.Fatal(err)
				}
				writer.Flush()
			}
		}
	}

	//creating a new writer, that will write into our csv file)
	writer := csv.NewWriter(file)
	writer.WriteAll(products)
	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
	writer.Flush()

	return true
}

// Function to Add Data Entries to
func writeKvData(option int, barcode string, valid bool) bool {

	if !valid {
		return true
	}

	if barcode == "end" {
		return false
	}

	err := bdb.View(func(txn *badger.Txn) error {

		txn = bdb.NewTransaction(false)
		_, err := txn.Get([]byte(barcode + "Name"))
		if err == badger.ErrKeyNotFound {
			name, category, description := getParams()
			// Initialize a Read-Write Transaction
			err = bdb.Update(func(txn *badger.Txn) error {

				// Create the Transaction, make it writable.
				txn = bdb.NewTransaction(true)
				defer txn.Discard()

				// Use the Transaction
				err := txn.Set([]byte(barcode+"Name"), []byte(name))
				check(err)
				err = txn.Set([]byte(barcode+"Category"), []byte(category))
				check(err)
				err = txn.Set([]byte(barcode+"Description"), []byte(description))
				check(err)
				err = txn.Commit()
				check(err)

				return err
			})
			check(err)
		} else {
			fmt.Println("The Barcode ", barcode, "already exists in this Database.")
			if option != 0 {
				sleep()
			}
		}

		return err
	})
	check(err)

	return true
}

func charLimiter(s string, limit int) string {
	//create a new reader, that is gonna read through s string
	reader := strings.NewReader(s)
	//create a buffer (slice of bytes), who's size gonna be limited
	buff := make([]byte, limit)
	// Fill the buff slice initially with spaces so we can TrimSpace the Strings to display them easier.
	for i := len(s); i < len(buff); i++ {
		buff[i] = 32
	}
	//using ReadAtLeast we gonna read (s) into buff until it has read at least minimum byte (limit)
	//it will read also futher, but buff is limited by (limit) and it will not take more characters than that
	n, _ := io.ReadAtLeast(reader, buff, limit)
	if n != 0 {
		if len(s) > limit {
			fmt.Printf("You wrote %d character. Only %d of them will be accepted\n", len(s), limit)
		}
		return string(buff)
	}
	return s

}

func addSQL(db *sql.DB, code string) {

	err := db.QueryRow("SELECT product_code FROM product_data WHERE product_code = ?", code).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {

			name, category, description := getParams()

			_, err := db.Exec("INSERT INTO product_data VALUES (?, ?, ?, ?)", strings.TrimSpace(code), strings.TrimSpace(name), strings.TrimSpace(category), strings.TrimSpace(description))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("You have successfully added an Item! You will be redirected to the main menu")
			sleep()
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("This Item has been stored already. You will be redirected to the main menu")
		sleep()
	}
}
