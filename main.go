package main

import (
	"fmt"
	"os"
	"urlshortener/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
