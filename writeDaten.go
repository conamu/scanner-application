package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func writeDaten() {

	var path = "testDatabase.csv"
	//open a file with flags: to append (O_Append) and to write(O_WRONLY)
	//FileMode (permission) - to append only
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//creating a two-demensional slice, where we gonna save input-sclices
	var products [][]string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Please, scan a barcode: ")
	scanner.Scan()
	barcode := scanner.Text()

	fmt.Print("Please, write a name: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Please, write a category: ")
	scanner.Scan()
	category := scanner.Text()

	fmt.Print("Please, write a description: ")
	scanner.Scan()
	description := scanner.Text()

	//creating a slice "product", which holds input as it's elements
	product := []string{barcode, name, category, description}

	//adding a slice "product" to two-dementtional slice "products"
	products = append(products, product)

	//creating a new writer, that will write to our csv file)
	writer := csv.NewWriter(file)
	writer.WriteAll(products)

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

}
