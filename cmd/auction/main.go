package main

import (
	"log"
	"main/bootstrapper"
)

func main() {
	err := bootstrapper.New().RunAPI()
	if err != nil {
		log.Fatal("failed to start")
	}
}
