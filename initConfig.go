package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("useKeyValueDB", true)
	viper.SetDefault("useFlatDB", false)
	viper.SetDefault("dbPath", "data/badger.db")
}
