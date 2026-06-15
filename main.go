/*
Copyright © 2026 Nitin Chouhan <developer.nitinchouhan@gmail.com>
*/
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nitinchouhan1/cloudctl/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cmd.Execute()
}
