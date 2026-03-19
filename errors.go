package main

import "fmt"

var hadError bool = false

// TODO: maybe an ErrorReporter interface to abstract
// how we report the error ?

// printError prints the error at line with msg
func printError(line uint32, msg string) {
	report(line, "", msg)
}

// report reports the error at line and flag the program to have error
func report(line uint32, where string, msg string) {
	fmt.Printf("[line %v] - error %v: %v\n", line, where, msg)
	hadError = true
}
