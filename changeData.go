package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/dgraph-io/badger/v2"
)

func chooseColumn() []string {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter or scan a code.")
	scanner.Scan()
	_, record := csvRead(scanner.Text(), "5")

	fmt.Printf(`Please choose which column you want to change.
	If you want to change the Name, press 1;
	If you want to change the Category, press 2;
	If you want to change the Description, press 3.
	I want to change: `)
	scanner.Scan()
	option, _ := strconv.Atoi(scanner.Text())

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
	default:
		fmt.Println("Invalid operation.")
		chooseColumn()
	}

	return record
}

func editKVEntry() {
	err := db.Update(func(txn *badger.Txn) error {
		txn = db.NewTransaction(true)
		barcode := getBarcode()

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
				continue
			}
		}

		txn.Commit()

		return err
	})
	check(err)
}
