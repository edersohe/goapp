package main

import "github.com/spf13/viper"

func main() {
	config()
	routes("")

	log.Fatal(e.Start(":" + viper.GetString("PORT")))
}
