package main

import "os"

const BIT_INLINE uint8 = 1 << 0
const BIT_HELP uint8 = 1 << 1

// Convert command flags to a bitmask
func GetFlagBitmask() uint8 {
	if len(os.Args) == 0 {
		return 0
	}

	var mask uint8 = 0
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-c":
			mask = mask | BIT_INLINE
		case "-h":
			mask = mask | BIT_HELP
		case "-ch":
			mask = mask | BIT_INLINE | BIT_HELP
		}
	}

	return mask
}

const HelpText = `Sequential: An interactive task queue for single-core organisms

Usage:

    sequential

The flags are:

-c
    Skips interactive mode and prints the current task
    directly to the command line

-h
    Show help text

Sequential opens an interactive CLI that allows the user to
see a single task without distraction, queue additional tasks,
and randomly choose a different task.`
