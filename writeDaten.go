package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	fmt.Print("Please, write a name (max. 150 character): ")
	scanner.Scan()
	name := scanner.Text()
	name = charLimiter(name, 150)

	fmt.Print("Please, write a category (max. 20 character): ")
	scanner.Scan()
	category := scanner.Text()
	category = charLimiter(category, 20)

	fmt.Print("Please, write a description (max. 500 character): ")
	scanner.Scan()
	description := scanner.Text()
	description = charLimiter(description, 500)

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

func charLimiter(s string, limit int) string {
	//create a new reader, that is gonna read through s string
	reader := strings.NewReader(s)
	//create a buffer, who's size gonna be limited
	buff := make([]byte, limit)
	//using ReadAtLeast we gonna read (s) to buff until it has read at least minimum byte (limit)
	n, _ := io.ReadAtLeast(reader, buff, limit)
	if n != 0 {
		if len(s) > limit {
			fmt.Printf("You wrote %d character. Only %d of them will be accepted\n", len(s), limit)
		}
		return string(buff)
	} else {
		return s
	}
}
