package main

import (
	"log"
	"main/bootstrapper"
)

func main() {
	err := bootstrapper.New().RunAPI()
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
