package main

import "github.com/spf13/viper"

func initConfig() {
	// Set configuration Defaults
	viper.SetDefault("useFlatDB", false)
	viper.SetDefault("useKeyValueDB", true)
	viper.SetDefault("activateRestApi", false)
	viper.SetDefault("dbPath", "data/badger.db")
	viper.SetDefault("flatPath", "data/database.csv")

	// If it doesnt exist, create a new config file with the default values.
	viper.SafeWriteConfigAs(".config.yaml")

	// Read the config File
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	check(err)

}
