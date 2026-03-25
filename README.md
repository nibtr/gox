# gox

gox is a simple tree-walk interpreter for educational purposes, mainly
referencing the "Crafting Interpreter" book but with my own implementation
to port into go.

The project's purpose is for me to learn more about go and how programming
languages work in general, so bugs and imperfect code are to be
expected. It's not a full-fledge product anyway.

The remaining sections are mainly my notes and understandings on the
topic of interpreters.

## Features

- Dynamic typing
- GC
- Types:
  - Booleans
  - Numbers: double precision floating point
  - Strings
  - nil
- Expressions:
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

<source_code> -> lexing -> <tokens> -> parsing -> <syntax_tree> -> static
analysis -> <intermediate_representation> -> optimize -> code gen
(native or bytecode)

### Asides (read more)

Intermediate representation (IR):
- Control flow graph
- Static single-assignment
- Continuation-passing style
- Three-address code (TAC or 3AC)

Optimization:
- Constant propagation
- Common subexpression elimination
- Loop invariant code motion
- Global value numbering
- Strength reduction
- Scalar replacement of aggregates
- Dead code elimination
- Loop unrolling

## Tree-walk interpreter

Basically, the interpreter will begin executing code right after parsing
it into an AST by traversing the syntax tree one branch and leaf at a
time and evaluate each node as it goes.

It is simple enough for me to implement.

## Compiler vs Interpreter

- Compiler: Translates one source to another (usually machine code or bytecode)
  without executing it.
- Interpreter: Takes in source code and executes it immediately. It runs
  programs "from source".

## Grammar

### Some basics

Backus-Naur form (BNF):

- Terminal: Tokens coming from the lexer (`if` or `1234`)
- Nonterminal: A named reference to another rule in the grammar

For the notation (simplified): 

```
rule_name -> <sequence_of_symbols> ;
```

### Expressions

- Literals: numbers, strings, booleans, nil
- Unary expressions: a prefix ! to perform a logical not, and - to negate a
number.
- Binary expressions: infix arithmetic ( + , - , * , / ) and logic (==,
!=, <, <=, >, >=)
- Parentheses: ( and )

E.g: `1 - (2 * 3) < 4 == false`

expression -> literal
              | unary
              | binary
              | grouping ;

literal    -> NUMBER | STRING | "true" | "false" | "nil" ;
grouping   -> "(" expression ")"
unary      -> ("-" | "!") expression ;
binary     -> expression operator expression;
operator   -> "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-" | "*" | "/" ;

**Note**: CAPITALIZE terminals that are a single lexeme whose text
representation may vary. NUMBER is any number literal, and STRING is
any string literal.

## License

See [LICENSE](./LICENSE) for more information.
