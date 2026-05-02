package main

import (
	"bufio"
	"errors"
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

	err = run(string(bytes))
	var le *LexerError
	var pe *ParseError
	var re *RuntimeError

	if errors.As(err, &le) {
		os.Exit(65)
	}
	if errors.As(err, &pe) {
		os.Exit(66)
	}
	if errors.As(err, &re) {
		os.Exit(67)
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
func run(source string) error {
	l := newLexer(source)
	tokens, err := l.scanTokens()
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	parser := NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	// fmt.Println(astPrinter{}.Print(expr))
	intrp := interpreter{}
	value, err := intrp.Intepret(expr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Println(value)

	return nil
}

func main() {
	args := os.Args[1:]
	l := len(args)

	if l > 1 {
		fmt.Println("Usage: ./bin/gox [script]")
	} else if l == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
