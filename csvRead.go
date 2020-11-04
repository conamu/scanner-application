package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func csvRead(code string) {
	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)

	// If the code matches an entry in the Databease, show the data. Else return an error.
	for index, record := range records {
		if record[0] == code {
			fmt.Println("Nr.: ", index, " == ", record[0], " == ", record[1], " == ", record[2], "\nDescription: ", record[3])
		} else if index + 2 > len(records) && record[0] != code {
			fmt.Println("This code is not stored in the system.")
		}
	}
}