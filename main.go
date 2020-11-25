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

	if viper.GetBool("activateRestApi") {
		fmt.Println("Rest API v 2.0 - Mux Routers")
		handleRequests()
	}
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	initMenus()

	mainMenu := menustyling.GetStoredMenu("main")

	for true {
		mainMenu.DisplayMenu()
		switch mainMenu.GetInputData() {
		case "1": // Get Data of one Entry
			code, valid := getBarcode()
			if viper.GetBool("useKeyValueDB") {

				readKV(code, valid)
				sleep()

			} else if viper.GetBool("useFlatDB") {
				csvRead(code, mainMenu.GetInputData(), valid)
			}
		case "2": // Edit one Entry based on Barcode
			if viper.GetBool("useKeyValueDB") {
				barcode, valid := getBarcode()
				editKVEntry(barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				deleteData("", chooseColumn(), true)
			}
		case "3": // Delete one Entry based on Barcode
			fmt.Println("WARNING! CODE SCANNED WILL BE PERMANENTLY ERASED FROM DATABASE!")
			barcode, valid := getBarcode()
			if viper.GetBool("useKeyValueDB") {
				deleteBadger(barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				deleteData(barcode, []string{}, valid)
			}
		case "4": // Add one Entry
			barcode, valid := getBarcode()
			if viper.GetBool("useKeyValueDB") {
				writeKvData(0, barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				writeData([]string{}, barcode, valid)
			}
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			if viper.GetBool("useKeyValueDB") {
				loop := true
				for loop {
					code, valid := getBarcode()
					loop, _ = readKV(code, valid)

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
					barcode, valid := getBarcode()
					loop = writeKvData(1, barcode, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				for loop {
					barcode, valid := getBarcode()
					loop = writeData([]string{}, barcode, valid)
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
