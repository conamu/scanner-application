package main

import "github.com/spf13/viper"

func initConfig() {
	viper.SetDefault("useKeyValueDB", false)
	viper.SetDefault("useFlatDB", true)
	viper.SetDefault("dbPath", "data/badger.db")
}
