package main

import (
	"encoding/csv"
	"os"
)

func deleteData(code string) {
	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR, 0755)
	defer file.Close()
	check(err)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)
	err = os.Remove("data/testDatabase.csv")
	check(err)

	nFile, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	defer nFile.Close()
	check(err)

	writer := csv.NewWriter(nFile)

	for _, record := range records {
		if record[0] != code {
			writer.Write(record)
		}
	}
	writer.Flush()
}