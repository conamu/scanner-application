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

	for index, record := range records {
		if record[0] == code {
			fmt.Println("Nr.: ", index, " == ", record[0], " == ", record[1], " == ", record[2], "\ndescription: ", record[3])
		}
	}
}