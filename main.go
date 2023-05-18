package main

import (
	"log"
)

func main() {
	r := SetupServer()
	if err := r.RunTLS(":443", "key.pem", "key.key"); err != nil {
		log.Printf("err: %v\n", err)
		return
	}
}
