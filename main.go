package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/per1Peteia/rfl/internal/cli"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	cli.InitCLI()
}
