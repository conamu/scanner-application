package main

import "github.com/conamu/cliutilsmodule/menustyling"

func initMenus() {
	menuText := make([][]string, 3)
	menuText[0] = []string{"Welcome!", "Please choose an option."}
	menuText[1] = []string{"1) Scan a Barcode", "2) Edit an Entry", "3) Delete an Entry", "4) Add an Entry", "5) Endless Scanning","6) Endless Adding", "q) Quit"}
	menuText[2] = []string{"Note: The endless modes can only", "be closed with a closing barcode or strg+c"}
	menustyling.CreateMenu(menuText, "=", 2, 1, true, true).StoreMenu("main")
}

