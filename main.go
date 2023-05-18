package main

import (
	"log"
)

func main() {
	r := SetupServer()
	if err := r.Run(":80"); err != nil {
		log.Printf("err: %v\n", err)
		return
	}
}
