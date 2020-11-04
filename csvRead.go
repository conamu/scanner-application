package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func csvRead(code string, option string) {
	file, err := os.OpenFile("data/testDatabase.csv", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	check(err)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	check(err)

	// If the code matches an entry in the Databease, show the data. Else return an error.
	for index, record := range records {
		if record[0] == code {
			fmt.Println("Nr.: ", index,
				" == ", record[0],
				" == ", record[1],
				" == ", record[2],
				"\nDescription: ", record[3],
				"\n========================================================")
		} else if code == "end" {
			log.Println("Scanned end/exit code, exiting!\nBye!")
			os.Exit(0)
		} else if index + 2 > len(records) && record[0] != code {
			fmt.Println("This code is not stored in the system.")
		}
	}

	if option != "5" && option != "6" {
		time.Sleep(time.Second*4)
	}
}