package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Abstraction for a simple newline-separated logfile.
type Logfile struct {
	fh    *os.File
	Lines []string
}

// Load a logfile's contents and split it into lines for
// the struct.
func LoadLogfile(filename string) *Logfile {
	fh := ensureFileExists(filename)

	reader := bufio.NewReader(fh)
	lines := make([]string, 0)
	var (
		err  error = nil
		line []byte
	)
	for err == nil {
		line, _, err = reader.ReadLine()
		if len(line) > 0 {
			lines = append(lines, string(line))
		}
	}

	return &Logfile{
		fh:    fh,
		Lines: lines,
	}
}

// Stringify the log lines and overwrite logfile.
func (l *Logfile) Save() {
	if l.fh != nil {
		l.fh.Truncate(0)
		_, err := l.fh.Write([]byte(strings.Join(l.Lines, "\n")))
		assertSuccess(err, fmt.Sprintf("Error writing to file %s", l.fh.Name()))
	}
}

// Close the file handler and nil the handle property.
func (l *Logfile) Close() {
	if l.fh != nil {
		err := l.fh.Close()
		assertSuccess(err, fmt.Sprintf("Cannot close file handle for %s", l.fh.Name()))
		l.fh = nil
	}
}

// Retrieve or create a file in the ~/.sequential directory. If it cannot
// be created or opened, exit the program.
func ensureFileExists(filename string) *os.File {
	home, err := os.UserHomeDir()
	assertSuccess(err, "User home directory was not found")

	path := fmt.Sprintf("%s%s%s", home, string(os.PathSeparator), ".sequential")
	err = os.MkdirAll(path, os.ModePerm)
	assertSuccess(err, "Unable to create home sequential directory")

	path = fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	assertSuccess(err, "Cannot create triage file")

	return f
}
