package main

import (
	"imshit/aircraftwar/db"
	"imshit/aircraftwar/utils"
	"log"
)

func main() {
	utils.LoadConfig()
	db.InitRedis()
	db.InitSql()
	r := SetupServer()
	if err := r.RunTLS(":443", "key.pem", "key.key"); err != nil {
		log.Printf("err: %v\n", err)
		return
	}
}
