package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"

	"github.com/conamu/cliutilsmodule/menustyling"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initDB() *badger.DB {
	initConfig()
	var rdb *badger.DB = nil
	if viper.GetBool("useKeyValueDB") {
		db, err := badger.Open(badger.DefaultOptions(viper.GetString("dbPath")))
		check(err)
		rdb = db
	}
	return rdb
}

var db *badger.DB = initDB()

func main() {
	defer db.Close()
	initMenus()

	scanner := bufio.NewScanner(os.Stdin)
	mainMenu := menustyling.GetStoredMenu("main")

	for true {
		mainMenu.DisplayMenu()
		switch mainMenu.GetInputData() {
		case "1": // Get Data of one Entry
			fmt.Println("Please enter or scan a code.")
			scanner.Scan()
			if viper.GetBool("useKeyValueDB") {

				barcode := scanner.Text()
				//readBadger(barcode)
				if checkItem(barcode) == true {
					name := readName(barcode)
					category := readCat(barcode)
					description := readDes(barcode)
					//in next 2 prints I am trying to make the ouptut looks like in flat-DB, but it doesn't really work
					//i just leave it to check later
					//fmt.Printf("====================================================\n%s == %s == %s\nDescription: %s\n", barcode, name, category, description)
					/* 	fmt.Println("====================================================\n",
					string(barcode),
					" == ", string(name), "hey", " == ", string(category), "hey",
					"\nDescription: ", string(description),
					"\n====================================================") */
					fmt.Printf("Barcode: %s\nName: %s\nCategory: %s\nDescription: %s\n", barcode, name, category, description)
					time.Sleep(time.Second * 4)
				}

			} else if viper.GetBool("useFlatDB") {
				csvRead(scanner.Text(), mainMenu.GetInputData())
			}
		case "2": // Edit one Entry based on Barcode
			if viper.GetBool("useKeyValueDB") {
				editKVEntry()
			} else if viper.GetBool("useFlatDB") {
				deleteData("", chooseColumn())
			}
		case "3": // Delete one Entry based on Barcode
			fmt.Println("WARNING! CODE SCANNED WILL BE PERMANENTLY ERASED FROM DATABASE!")
			fmt.Println("Please enter or scan a code.")
			scanner.Scan()
			if viper.GetBool("useKeyValueDB") {
				deleteBadger(scanner.Text())
			} else if viper.GetBool("useFlatDB") {
				deleteData(scanner.Text(), []string{})
			}
		case "4": // Add one Entry
			if viper.GetBool("useKeyValueDB") {
				writeKvData(0)
			} else if viper.GetBool("useFlatDB") {
				writeDaten([]string{})
			}
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			if viper.GetBool("useKeyValueDB") {
				// function for KeyValue DB
			} else if viper.GetBool("useFlatDB") {
				loop := true
				for loop {
					scanner.Scan()
					loop, _ = csvRead(scanner.Text(), mainMenu.GetInputData())
				}
			}
		case "6": // Add endless entries, terminate with strg+c or "end" code
			loop := true
			if viper.GetBool("useKeyValueDB") {
				for loop {
					loop = writeKvData(1)
				}
			} else if viper.GetBool("useFlatDB") {
				for loop {
					loop = writeDaten([]string{})
				}
			}
		case "q": // Quit programm
			log.Println("pressed exit, programm Exiting.\nBye!")
			os.Exit(0)
		default:
			continue
		}
	}
}
