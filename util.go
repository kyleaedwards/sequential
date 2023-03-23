package main

import (
	"fmt"
	"os"
)

// Fatal an unexpected error and print a message to the user.
func assertSuccess(err error, str string) {
	if err != nil {
		fmt.Println(str)
		os.Exit(1)
	}
}
