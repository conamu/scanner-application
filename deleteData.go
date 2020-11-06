package main

import (
	"encoding/csv"
	"os"
)

func deleteData(code string, newRecord []string) {
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

	if code == "" {
		var indexToEdit int
		for index, record := range records {
			if record[0] == newRecord[0] {
				indexToEdit = index
			}
		}

		for i := 1; i < len(records[indexToEdit]); i++ {
			records[indexToEdit][i] = newRecord[i]
		}

		writer.WriteAll(records)

	} else {
		for _, record := range records {
			if record[0] != code {
				writer.Write(record)
			}
		}
	}

	writer.Flush()
}