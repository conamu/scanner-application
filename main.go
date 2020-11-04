package main

import (
	"bufio"
	"fmt"
	"github.com/conamu/cliutilsmodule/menustyling"
	"log"
	"os"
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
	mainM := menustyling.GetStoredMenu("main")
	mainM.DisplayMenu()

	switch mainM.GetInputData() {
	case "1": // Get Data of one Entry
		fmt.Println("Please enter or scan a code.")
		scanner.Scan()
		csvRead(scanner.Text())
	case "2": // Edit one Entry based on Barcode
	case "3": // Delete one Entry based on Barcode
	case "4": // Add one Entry
	case "5": // Get data from endless codes, terminate with strg+c or "end" code
	case "6": // Add endless entries, terminate with strg+c or "end" code
	case "q": // Quit programm
		log.Println("pressed exit, programm Exiting.\nBye!")
		os.Exit(0)
	default:
		main()
	}
}