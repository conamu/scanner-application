package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/conamu/cliutilsmodule/menustyling"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {

	initMenus()
	scanner := bufio.NewScanner(os.Stdin)
	mainMenu := menustyling.GetStoredMenu("main")

	for true {
		mainMenu.DisplayMenu()
		switch mainMenu.GetInputData() {
		case "1": // Get Data of one Entry
			fmt.Println("Please enter or scan a code.")
			scanner.Scan()
			csvRead(scanner.Text(), mainMenu.GetInputData())
		case "2": // Edit one Entry based on Barcode
			writeDaten(chooseColumn())
		case "3": // Delete one Entry based on Barcode
			fmt.Println("WARNING! CODE SCANNED WILL BE PERMANENTLY ERASED FROM DATABASE!")
			fmt.Println("Please enter or scan a code.")
			scanner.Scan()
			deleteData(scanner.Text())
		case "4": // Add one Entry
			var empty []string
			writeDaten(empty)
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			for true {
				scanner.Scan()
				csvRead(scanner.Text(), mainMenu.GetInputData())
			}
		case "6": // Add endless entries, terminate with strg+c or "end" code
		case "q": // Quit programm
			log.Println("pressed exit, programm Exiting.\nBye!")
			os.Exit(0)
		default:
			continue
		}
	}
}
