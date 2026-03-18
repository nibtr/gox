package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	args := os.Args[1:]
	l := len(args)

	if l > 1 {
		fmt.Println("Usage: go run ./main.go [script]")
	} else if l == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

// runFile runs the interpreter for the file at `path`
func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read file %s: %w", path, err))
	}

	run(string(bytes))
}

// runPrompt runs the interpreter for the current prompt from user
func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			// EOF or error
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			}
			break
		}

		run(scanner.Text())
	}
}

// run executes the interpreter for a source
func run(source string) {
	fmt.Println(source)
}
