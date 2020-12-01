package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/conamu/cliutilsmodule/menustyling"
	"github.com/dgraph-io/badger/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
		openSQL()
	}
}

func openSQL() *sql.DB {
	source := viper.GetString("mysqlUser") + ":" + viper.GetString("mysqlPassword") +
		"@tcp(" + viper.GetString("mysqlServerAddress") + ":" + viper.GetString("mysqlServerPort") +
		")/" + viper.GetString("mysqlDatabaseName")

	db, err := sql.Open("mysql", source)
	check(err)
	mdb = db
	_, err = mdb.Exec("create table IF NOT EXISTS product_data(product_code varchar(10), product_name varchar(150), product_category varchar(20), product_description varchar(500))")
	return mdb
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
			barcode, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				readKV(barcode, valid)
				sleep()
			} else if viper.GetBool("useFlatDB") {
				csvRead(barcode, mainMenu.GetInputData(), valid)
			}
		case "2": // Edit one Entry based on Barcode
			barcode, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				editKVEntry(barcode, valid)
			} else if viper.GetBool("useFlatDB") {
				deleteData("", chooseColumn(), true)
			} else if viper.GetBool("useMysqlDB") {
				updateSQL(openSQL(), barcode, valid)
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
			} else if viper.GetBool("useMysqlDB") {
				addSQL(openSQL(), barcode, valid)
				sleep()
			}
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			code, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				loop := true
				for loop {
					loop, _, _ = readKV(code, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				loop := true
				for loop {
					loop, _, _ = csvRead(code, mainMenu.GetInputData(), valid)
					if !valid {
						loop = true
					}
				}
			}
		case "6": // Add endless entries, terminate with strg+c or "end" code
			loop := true
			barcode, valid := validateBarcode(getBarcode())
			if viper.GetBool("useKeyValueDB") {
				for loop {
					loop = writeKvData(1, barcode, valid)
				}
			} else if viper.GetBool("useFlatDB") {
				for loop {
					loop = writeData([]string{}, barcode, valid)
				}
			} else if viper.GetBool("usemysqldb") {
				for loop {
					loop = addSQL(openSQL(), barcode, valid)
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
