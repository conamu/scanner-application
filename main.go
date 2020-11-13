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
func sleep() {
	time.Sleep(time.Second * 3)
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
var scanner = bufio.NewScanner(os.Stdin)

func main() {

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	initMenus()
  
	mainMenu := menustyling.GetStoredMenu("main")

	for true {
		mainMenu.DisplayMenu()
		switch mainMenu.GetInputData() {
		case "1": // Get Data of one Entry
			fmt.Println("Please enter or scan a code.")
			code, valid := getBarcode()
			if viper.GetBool("useKeyValueDB") {

				barcode := scanner.Text()

				//if checkItem(barcode) {
				/* name := readName(barcode)
				category := readCategory(barcode)
				description := readDescription(barcode) */

				//in next 2 prints I am trying to make the ouptut looks like in flat-DB, but it doesn't really work
				//i just leave it to check later
				//WHY IT WORKS AS IT DOES
				//fmt.Printf("====================================================\n%s == %s == %s\nDescription: %s\n", barcode, name, category, description)
				/* 	fmt.Println("====================================================\n",
				string(barcode),
				" == ", string(name), "hey", " == ", string(category), "hey",
				"\nDescription: ", string(description),
				"\n====================================================") */

				readKV(barcode)
				sleep()

			} else if viper.GetBool("useFlatDB") {
				csvRead(code, mainMenu.GetInputData(), valid)
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
				loop := true
				for loop {
					fmt.Println("Please enter or scan a code.")
					scanner.Scan()
					code := scanner.Text()
					loop = readKV(code)

				}

			} else if viper.GetBool("useFlatDB") {
				loop := true
				for loop {
					code, valid := getBarcode()
					loop, _, _ = csvRead(code, mainMenu.GetInputData(), valid)
					if !valid {
						loop = true
					}
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
			if viper.GetBool("useKeyValueDB") {
				db.Close()
			}
			os.Exit(0)
		default:
			continue
		}
	}
}
