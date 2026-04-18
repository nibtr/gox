package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
)

// runFile runs the interpreter for the file at `path`
func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read file %s: %w", path, err))
	}

	run(string(bytes))
	if hadError {
		os.Exit(1)
	}
}

// runPrompt runs the interpreter for the current prompt from user
func runPrompt() {
	// TODO: maybe use github.com/chzyer/readline to handle arrow keys?
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
		hadError = false // reset the error flag
	}
}

// run executes the interpreter for a source
func run(source string) {
	l := newLexer(source)
	tokens := l.scanTokens()
	parser := newParser(tokens)
	expr, err := parser.Parse()

	if err != nil {
		return
	}

	fmt.Println(astPrinter{}.print(expr))

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }
}

func main() {
	args := os.Args[1:]
	l := len(args)

	if l > 1 {
		fmt.Println("Usage: ./goitr [script]")
	} else if l == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
