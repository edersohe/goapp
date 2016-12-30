package main

import "github.com/spf13/viper"

func main() {
	config()
	routes("")

	e.Logger.Fatal(e.Start(":" + viper.GetString("PORT")))
}
