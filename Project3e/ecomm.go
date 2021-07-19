//main.go
package main

import (
	"fmt"

	"Project3e/Config"
	"Project3e/Models"
	"Project3e/Routes"

	"github.com/jinzhu/gorm"
)

var err error

func main() {
	Config.DB, err = gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer Config.DB.Close()
	Config.DB.AutoMigrate(&Models.Customer{})
	Config.DB.AutoMigrate(&Models.Product{})
	Config.DB.AutoMigrate(&Models.Cart{})
	Config.DB.AutoMigrate(&Models.Payment{})
	r := Routes.SetupRouter()
	//running
	r.Run()
}
