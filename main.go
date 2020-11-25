package main

import (
	"bufio"
	"fmt"
	"github.com/conamu/cliutilsmodule/menustyling"
	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
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
	if viper.GetBool("apiEndpointMode") {
		requestHandler()
	} else if !viper.GetBool("apiEndpointMode") {
		completeMode()
	}
}

func completeMode() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	initMenus()

	mainMenu := menustyling.GetStoredMenu("main")

	for true {
		mainMenu.DisplayMenu()
		switch mainMenu.GetInputData() {
		case "1": // Get Data of one Entry
			code, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				readKV(code, valid)
				sleep()
			} else if viper.GetBool("useFlatDB") {
				_, result, _ := csvRead(code, mainMenu.GetInputData(), valid)
				itemDisplay(result[1], result[2], result[3])
			}
		case "2": // Edit one Entry based on Barcode

			if viper.GetBool("useKeyValueDB") {
				barcode, valid := validateBarcode(getBarcode())
				editKVEntry(barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				deleteData("", chooseColumn(), true)
			}
		case "3": // Delete one Entry based on Barcode
			fmt.Println("WARNING! CODE SCANNED WILL BE PERMANENTLY ERASED FROM DATABASE!")
			barcode, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				deleteBadger(barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				deleteData(barcode, []string{}, valid)
			}
		case "4": // Add one Entry
			barcode, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				writeKvData(0, barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				writeData([]string{}, barcode, valid)
			}
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			if viper.GetBool("useKeyValueDB") {
				loop := true
				for loop {
					code, valid := validateBarcode(getBarcode())
					loop = readKV(code, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				loop := true
				for loop {
					code, valid := validateBarcode(getBarcode())
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
					barcode, valid := validateBarcode(getBarcode())
					loop = writeKvData(1, barcode, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				for loop {
					barcode, valid := validateBarcode(getBarcode())
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
