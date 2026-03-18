package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
)

type program struct {
	hadError bool
}

// runFile runs the interpreter for the file at `path`
func (p *program) runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read file %s: %w", path, err))
	}

	p.run(string(bytes))
	if p.hadError {
		os.Exit(65)
	}
}

// runPrompt runs the interpreter for the current prompt from user
func (p *program) runPrompt() {
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

		p.run(scanner.Text())
		p.hadError = false // reset the error flag
	}
}

// run executes the interpreter for a source
func (p *program) run(source string) {
	l := newLexer(source)
	tokens := l.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

// printError prints the error at line with msg
func (p *program) printError(line uint32, msg string) {
	p.report(line, "", msg)
}

// TODO: maybe an ErrorReporter interface to abstract
// how we report the error ?

// report reports the error at line and flag the program to have error
func (p *program) report(line uint32, where string, msg string) {
	fmt.Printf("[line %v] - error %v: %v", line, where, msg)
	p.hadError = true
}

func main() {
	args := os.Args[1:]
	l := len(args)

	if l > 1 {
		fmt.Println("Usage: ./goitr [script]")
	} else if l == 1 {
		p := &program{}
		p.runFile(args[0])
	} else {
		p := &program{}
		p.runPrompt()
	}
}
