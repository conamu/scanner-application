package main

import (
	"fmt"
	"log"
	"strings"
)

func getBarcode() string {
	fmt.Print("Scan or enter a Barcode: ")
	scanner.Scan()
	return scanner.Text()
}

func validateBarcode(inputCode string) (string, bool) {

	valid := false
	code := ""

	if strings.HasPrefix(inputCode, "H24") && len(inputCode) < 10 {
		code = inputCode
		valid = true
	} else if inputCode == "end" {
		code = inputCode
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
