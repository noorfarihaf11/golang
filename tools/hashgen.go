package main

import (
	"fmt"
	"log"

	"github.com/noorfarihaf11/clean-arc/utils"
)

func main() {
	hash, err := utils.HashPassword("123456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash password:", hash)
}
