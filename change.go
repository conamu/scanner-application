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
	record := csvRead(scanner.Text(), "5")

	fmt.Printf(`Please, choose which column you want to change.
	If you want to change Name press 1;
	If you want to change Category, press 2;
	If you want to change Description, press 3.
	I want to change: `)
	scanner.Scan()
	option, _ := strconv.Atoi(scanner.Text())

	switch option {
	case 1:
		fmt.Println("Please, write new Name: ")
		scanner.Scan()
		newName := scanner.Text()
		record[1] = charLimiter(newName, 150)
	case 2:
		fmt.Println("Please, write new Category: ")
		scanner.Scan()
		newCategory := scanner.Text()
		record[2] = charLimiter(newCategory, 20)
	case 3:
		fmt.Println("Please, write new Description: ")
		scanner.Scan()
		newDescr := scanner.Text()
		record[3] = charLimiter(newDescr, 500)
	default:
		fmt.Println("Fail")

	}
	/* 	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	   	defer file.Close()
	   	check(err)
	   	reader := csv.NewReader(file)
	   	records, err := reader.ReadAll()
	   	check(err)

	   	fmt.Println(records)

	   	writer := csv.NewWriter(file)
	   	writer.Write(record)

	   	if err := writer.Error(); err != nil {
	   		log.Fatal(err)
	   	}
	   	fmt.Println("========================================")
	   	fmt.Println(records) */
	fmt.Println(record)
	return record
}
