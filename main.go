package main

import (
	"flag"
	"fmt"
	"go-rgw/connection"
	"go-rgw/gc"
	"go-rgw/router"
)

var defaultConfig = "./application.yml"

func main() {
	configFile := flag.String("config", defaultConfig, "configuration filename")
	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Println(err)
	}
	mysql := connection.NewMySQL(config.Database.Username, config.Database.Password, config.Database.Address, config.Database.Name,
		"utf8mb4")
	err = mysql.Init()
	if err != nil {
		fmt.Println(err)
	}
	connection.InitMySQLManager(mysql)
	ceph, err := connection.NewCeph()
	if err != nil {
		fmt.Println(err)
	}
	err = ceph.InitDefault()
	if err != nil {
		fmt.Println(err)
	}
	connection.InitCephManager(ceph)
	gc.Init()
	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
