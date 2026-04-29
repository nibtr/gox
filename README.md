# gox

gox is a simple tree-walk interpreter for educational purposes, mainly
referencing the "Crafting Interpreter" book but with my own implementation
to port into go.

The project's purpose is for me to learn more about go and how programming
languages work in general, so bugs and imperfect code are to be
expected. It's not a full-fledge product anyway.

**Note**: Still a work in progress.

## Tree-walk interpreter

Basically, the interpreter will begin executing code right after parsing
it into an AST by traversing the syntax tree one branch and leaf at a
time and evaluate each node as it goes. It is simple enough for me to implement.

## Features

- Dynamic typing
- GC
- Types:
  - Booleans
  - Numbers: double precision floating point
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
- Standard lib

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
