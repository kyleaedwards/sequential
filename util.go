package main

import (
	"fmt"
	"os"
)

// Fatal an unexpected error and print a message to the user.
func AssertSuccess(err error, str string) {
	if err != nil {
		fmt.Println(str)
		os.Exit(1)
	}
}

// Retrieve or create a file in the ~/.sequential directory. If it cannot
// be created or opened, exit the program.
func EnsureFileExists(filename string) *os.File {
	home, err := os.UserHomeDir()
	AssertSuccess(err, "User home directory was not found")

	path := fmt.Sprintf("%s%s%s", home, string(os.PathSeparator), ".sequential")
	err = os.MkdirAll(path, os.ModePerm)
	AssertSuccess(err, "Unable to create home sequential directory")

	path = fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	AssertSuccess(err, fmt.Sprintf("Cannot create file %s", filename))

	return f
}
