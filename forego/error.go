package main

import (
	"fmt"
	"os"
)

func handleError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
