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

			deleteData("" ,chooseColumn())

		case "3": // Delete one Entry based on Barcode
			fmt.Println("WARNING! CODE SCANNED WILL BE PERMANENTLY ERASED FROM DATABASE!")
			fmt.Println("Please enter or scan a code.")
			scanner.Scan()
			deleteData(scanner.Text(), []string{})
		case "4": // Add one Entry
			writeDaten([]string{})
		case "5": // Get data from endless codes, terminate with strg+c or "end" code
			loop := true
			for loop {
				scanner.Scan()
				loop, _ = csvRead(scanner.Text(), mainMenu.GetInputData())
			}
		case "6": // Add endless entries, terminate with strg+c or "end" code
			loop := true
			for loop {
				loop = writeDaten([]string{})
			}
		case "q": // Quit programm
			log.Println("pressed exit, programm Exiting.\nBye!")
			os.Exit(0)
		default:
			continue
		}
	}
}
