package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func chooseColumn() []string {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter or scan a code.")
	scanner.Scan()
	_, record := csvRead(scanner.Text(), "5")

	fmt.Printf(`Please choose which column you want to change.
	If you want to change the Name, press 1;
	If you want to change the Category, press 2;
	If you want to change the Description, press 3.
	I want to change: `)
	scanner.Scan()
	option, _ := strconv.Atoi(scanner.Text())

	switch option {
	case 1:
		fmt.Println("Enter a New Product Name: ")
		scanner.Scan()
		newName := scanner.Text()
		record[1] = charLimiter(newName, 150)
	case 2:
		fmt.Println("Enter a new Category:  ")
		scanner.Scan()
		newCategory := scanner.Text()
		record[2] = charLimiter(newCategory, 20)
	case 3:
		fmt.Println("Enter a new Description: ")
		scanner.Scan()
		newDescr := scanner.Text()
		record[3] = charLimiter(newDescr, 500)
	default:
		fmt.Println("Invalid operation.")
		chooseColumn()
	}

	return record
}

func editKVEntry() {
	
}
