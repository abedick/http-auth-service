package main

import "fmt"

var c serverConfig

func main() {
	path := "./conf/auth.toml"
	parseConfig(path)
	c.ConfigPath = path

	err := c.Logger.Init()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("cannot start service without logging capabilities")
		return
	}
	c.Logger.Message("success", "service started")

	initServer()
}
