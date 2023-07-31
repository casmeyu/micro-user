package main

import (
	"fmt"
	"log"
	"os"

	"github.com/casmeyu/micro-user/configuration"
	structs "github.com/casmeyu/micro-user/structs"
)

// Setting config as global variable
var Config structs.Config

func executeAppSetup() {
	err := configuration.LoadConfig(&Config)
	if err != nil {
		log.Println("Error while setting up config", err.Error())
		os.Exit(2)
	}
}

func main() {
	executeAppSetup()
	fmt.Println("Configuration loaded")
	fmt.Println(Config)
}
