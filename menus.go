package main

import (
	"github.com/conamu/cliutilsmodule/menustyling"
	"strings"
)

func initMenus() {
	menuText := make([][]string, 3)
	menuText[0] = []string{"Welcome to the Scanning Application!", "Please choose an option."}
	menuText[1] = []string{"1) Scan a Barcode", "2) Edit an Entry", "3) Delete an Entry", "4) Add an Entry", "5) Endless Scanning","6) Endless Adding", "q) Quit"}
	menuText[2] = []string{"Note: The endless modes can only be closed", "by typing or scanning \"end\""}
	menustyling.CreateMenu(menuText, "=", 2, 1, true, true).StoreMenu("main")
}

func itemDisplay(name, category, description string) {
	menuText := make([][]string, 1)
	menuText[0] = []string{"Item: " + strings.TrimSpace(name) + "  == " + " Category: " + strings.TrimSpace(category), "Item Description: " + strings.TrimSpace(description)}
	menustyling.CreateMenu(menuText, "=", 1, 1, false, false).DisplayMenu()
}

func itemEditMenu() string {
	menuText := make([][]string, 2)
	menuText[0] = []string{"Please choose which value you want to edit."}
	menuText[1] = []string{"1) Name", "2) Category", "3) Description", "4) Abort"}
	menu := menustyling.CreateMenu(menuText, "=", 1, 1, true, true)
	menu.DisplayMenu()
	return menu.GetInputData()
}

