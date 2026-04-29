# gox

gox is a simple tree-walk interpreter written in Go.

The project is inspired by the book *Crafting Interpreters* but
with my own implementation to learn Go and understand how programming
languages work internally.

This is a learning project and implementation details will evolve over time.

**Status**: Work in progress.

## Tree-walk interpreter

The interpreter executes code by first parsing source code into an AST
(Abstract Syntax Tree), then recursively evaluating each node. It is
simple enough to learn.

## Features

- Dynamic typing
- GC
- Types:
  - Booleans
  - Numbers: float64
  - Strings
  - nil
- Expressions:
  - Ternary
  - Arithmetic
  - Comparison
  - Logical operators
- Statements
- Variables
- Control flow
- Functions: first class -> Closures
- Classes
- Standard lib (planned)

## Process

<source_code> -> lexing -> <tokens> -> parsing -> <syntax_tree>(currently here) -> static
analysis -> <intermediate_representation> -> optimize -> code gen
(native or bytecode)

## Run

You can build to binary first by `go build -o bin/gox`, or can just run with `go
run ./...`. I will refactor to packages later.

```
go run ./... # run the repl

> 1 * (2 + 3) # 5

# or with files

go run ./... <file.gox> # still work in progress
```

## Grammar

See [GRAMMAR.md](/GRAMMAR.md) for more details.

## License

See [LICENSE](./LICENSE) for more information.
