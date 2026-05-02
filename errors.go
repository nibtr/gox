package main

import "fmt"

// TODO: maybe refactor this whole file lol

var hadError bool = false
var hadRuntimeError bool = false

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

func runtimeError(e *RuntimeError) {
	fmt.Printf("%v\n[line %v]", e.message, e.tok.line)
	hadRuntimeError = true
}
