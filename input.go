package main

import (
	"fmt"
	"log"
	"strings"
)

func getBarcode() (string, bool) {
	fmt.Print("Scan a Barcode: ")
	scanner.Scan()
	scannedBarcode := scanner.Text()

	valid := false
	code := ""

	if strings.HasPrefix(scannedBarcode, "H24") && len(scannedBarcode) < 10 {
		code = scannedBarcode
		valid = true
	} else if scannedBarcode == "end" {
		code = scannedBarcode
		valid = true
	} else {
		log.Println("Please Scan a valid Barcode.")
	}
	return code, valid
}

func getParams() (name, category, description string) {

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

	return name, category, description
}
