package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"io"
	"log"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func getParams() (barcode, name, category, description string) {

	fmt.Print("Scan a Barcode: ")
	scanner.Scan()
	barcode = scanner.Text()

	fmt.Print("Whats the Product Name? (max. 150 characters): ")
	scanner.Scan()
	name = scanner.Text()
	name = charLimiter(name, 150)

	fmt.Print("Whats the Product Category? (max. 20 characters): ")
	scanner.Scan()
	category = scanner.Text()
	category = charLimiter(category, 20)

	fmt.Print("Whats the Product Description? (max. 500 characters): ")
	scanner.Scan()
	description = scanner.Text()
	description = charLimiter(description, 500)

	return barcode, name, category, description
}

func writeDaten(data []string) bool {


	var path = "data/testDatabase.csv"
	//open a file with flags: to append (O_Append) and to write(O_WRONLY)
	//FileMode (permission) - to append only
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)

	//creating a two-dimensional slice, where we save input-slices
	var products [][]string

	barcode := ""

	if len(data) == 0 {

		barcode, name, category, description := getParams()
		if barcode == "end" {
			return false
		}

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

func writeKvData() {

	// barcode, name, category, description := getParams()



	// Initialize a Read-Write Transaction
	err := db.Update(func(txn *badger.Txn) error {
		// Create the Transaction, make it writable.
		txn = db.NewTransaction(true)
		defer txn.Discard()

		// Use the Transaction
		err := txn.Set([]byte("test"), []byte("This is a test Value"))
		check(err)
		err = txn.Set([]byte("test"), []byte("This is another test value."))
		check(err)
		err = txn.Commit()
		check(err)

		return err
	})
	check(err)

	db.View(func(txn *badger.Txn) error {

		txn = db.NewTransaction(false)
		item, _ := txn.Get([]byte("test"))

		value, _ := item.ValueCopy(nil)
		fmt.Println(string(value))

		return nil
	})


}

func charLimiter(s string, limit int) string {
	//create a new reader, that is gonna read through s string
	reader := strings.NewReader(s)
	//create a buffer (slice of bytes), who's size gonna be limited
	buff := make([]byte, limit)
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
