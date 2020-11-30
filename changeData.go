package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/dgraph-io/badger/v2"
)

func chooseColumn() []string {

	code, valid := validateBarcode(getBarcode())
	if !valid {
		return nil
	}
	_, record, err := csvRead(code, "5", valid)
	if err != nil {
		return nil
	}

	option, _ := strconv.Atoi(itemEditMenu())

	switch option {
	case 1:
		fmt.Println("Enter a New Product Name: ")
		scanner.Scan()
		newName := scanner.Text()
		record[1] = charLimiter(newName, 150)
	case 2:
		fmt.Println("Enter a new Category:  ")
		scanner.Scan()
		newCategory := scanner.Text()
		record[2] = charLimiter(newCategory, 20)
	case 3:
		fmt.Println("Enter a new Description: ")
		scanner.Scan()
		newDescr := scanner.Text()
		record[3] = charLimiter(newDescr, 500)
	case 4:
		return nil
	default:
		fmt.Println("Invalid operation.")
		chooseColumn()
	}

	return record
}

func editKVEntry(barcode string, valid bool) {

	if !valid {
		return
	}

	err := bdb.Update(func(txn *badger.Txn) error {
		txn = bdb.NewTransaction(true)
		if !valid {
			return errors.New("CODE NOT VALID")
		}

		nameVal, err := txn.Get([]byte(barcode + "Name"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("This Barcode does not Exist yet.\nPlease use the Add functionality to add it.")
			sleep()
			return nil
		}
		categoryVal, err := txn.Get([]byte(barcode + "Category"))
		descriptionVal, err := txn.Get([]byte(barcode + "Description"))
		check(err)

		nameOld, err := nameVal.ValueCopy(nil)
		categoryOld, err := categoryVal.ValueCopy(nil)
		descriptionOld, err := descriptionVal.ValueCopy(nil)
		check(err)

		itemDisplay(string(nameOld), string(categoryOld), string(descriptionOld))
		option := itemEditMenu()
		m := true
		for m {
			switch option {
			case "1":
				fmt.Println("Please enter a new Item Name: ")
				scanner.Scan()
				name := scanner.Text()
				err := txn.Set([]byte(barcode+"Name"), []byte(name))
				check(err)
				m = false
			case "2":
				fmt.Println("Please enter a new Item Category: ")
				scanner.Scan()
				category := scanner.Text()
				err := txn.Set([]byte(barcode+"Category"), []byte(category))
				check(err)
				m = false
			case "3":
				fmt.Println("Please enter a new Item Description: ")
				scanner.Scan()
				description := scanner.Text()
				err := txn.Set([]byte(barcode+"Description"), []byte(description))
				check(err)
				m = false
			case "4":
				return nil
			default:
				fmt.Println("You have choosen an anvalid option! You will be redirected to the main menu")
			}
		}

		txn.Commit()

		return err
	})
	check(err)
}

func updateSQL(db *sql.DB, barcode string, valid bool) {
	if barcode == "end" {
		return
	}

	if !valid {
		return
	}
	var name, category, description string
	err := db.QueryRow("SELECT product_name, product_category, product_description FROM product_data WHERE product_code = ?", barcode).Scan(&name, &category, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("The Barcode %s does not exist yet.\nPlease use the Add functionality to add it.\n", barcode)
			sleep()
			return
		}
		fmt.Println(name, category, description)
		log.Fatal(err)
	}

	itemDisplay(name, category, description)
	options := itemEditMenu()

	switch options {
	case "1":
		//Name
		fmt.Println("Please enter a new Item Name: ")
		scanner.Scan()
		name := scanner.Text()
		name = charLimiter(name, 150)
		_, err := db.Exec("UPDATE product_data SET product_name = ? WHERE product_code = ?", name, barcode)
		check(err)
	case "2":
		//Category
		fmt.Println("Please enter a new Item Category: ")
		scanner.Scan()
		category := scanner.Text()
		category = charLimiter(category, 20)
		_, err := db.Exec("UPDATE product_data SET product_category = ? WHERE product_code = ?", category, barcode)
		check(err)
	case "3":
		//Description
		fmt.Println("Please enter a new Item Description: ")
		scanner.Scan()
		description := scanner.Text()
		description = charLimiter(description, 500)
		_, err := db.Exec("UPDATE product_data SET product_description = ? WHERE product_code = ?", description, barcode)
		check(err)
	case "4":
		break
	default:
		fmt.Println("You have choosen an anvalid option! You will be redirected to the main menu")

	}
}
