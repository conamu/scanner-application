package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/dgraph-io/badger/v2"
	"github.com/spf13/viper"
	"github.com/conamu/cliutilsmodule/menustyling"
	_ "github.com/go-sql-driver/mysql"
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

// Use these global db objects in application.
var bdb *badger.DB
var mdb *sql.DB

func initDB() {
	if viper.GetBool("useKeyValueDB") {
		db, err := badger.Open(badger.DefaultOptions(viper.GetString("dbPath")))
		check(err)
		bdb = db
	// Construct a string with all necessary details about the mysql account and server to create a db object.
	// If it doesnt exist, create a table named product_data.
	} else if viper.GetBool("useMysqlDB") {
		source := viper.GetString("mysqlUser") + ":" + viper.GetString("mysqlPassword") +
			"@tcp(" + viper.GetString("mysqlServerAddress") + ":" + viper.GetString("mysqlServerPort") +
			")/" + viper.GetString("mysqlDatabaseName")

		fmt.Println(source)

		db, err := sql.Open("mysql", source)
		check(err)
		mdb = db
		res, err := mdb.Exec("create table IF NOT EXISTS product_data(product_code varchar(10), product_name varchar(150), product_category varchar(20), product_description varchar(200))")
		fmt.Println(res)
	}
}
var scanner = bufio.NewScanner(os.Stdin)

func main() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}
	initConfig()
	if !viper.GetBool("apiEndpointMode") {
		initDB()
	}
	if viper.GetBool("apiEndpointMode") {
		requestHandler()
	} else if !viper.GetBool("apiEndpointMode") {
		completeMode()
	}
}

func completeMode() {

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
				csvRead(code, mainMenu.GetInputData(), valid)
			} else if viper.GetBool("useMysqlDB") {
				_, record, err := readSql(code, valid)
				if errors.Is(err, notFound) {
					fmt.Println("This code is not stored in the Database")
				} else {
					itemDisplay(record[1], record[2], record[3])
					sleep()
				}
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
			} else if viper.GetBool("useMysqlDB") {
				deleteSql(barcode, valid)
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
					loop, _, _ = readKV(code, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				loop := true
				var record []string
				for loop {
					code, valid := validateBarcode(getBarcode())
					loop, record, _ = csvRead(code, mainMenu.GetInputData(), valid)
					if !valid {
						loop = true
					}
					itemDisplay(record[1], record[2], record[3])
				}
			} else if viper.GetBool("useMysqlDB") {
				loop := true
				var record []string
				var err error
				for loop {
					code, valid := validateBarcode(getBarcode())
					loop, record, err = readSql(code, valid)
					if errors.Is(err, notFound) {
						fmt.Println("This code is not stored in the Database")
					} else if loop {
						itemDisplay(record[1], record[2], record[3])
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
				bdb.Close()
			} else if viper.GetBool("useMysqlDB") {
				mdb.Close()
			}
			os.Exit(0)
		default:
			continue
		}
	}
}
