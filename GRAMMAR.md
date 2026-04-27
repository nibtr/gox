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

```
expression -> literal
              | unary
              | binary
              | grouping ;

literal    -> NUMBER | STRING | "true" | "false" | "nil" ;
grouping   -> "(" expression ")"
unary      -> ("-" | "!") expression ;
binary     -> expression operator expression;
operator   -> "==" | "!=" | "<" | "<=" | ">" | ">=" | "+" | "-" | "*" | "/" ;
```

**Note**: CAPITALIZE terminals that are a single lexeme whose text
representation may vary. NUMBER is any number literal, and STRING is
any string literal.

### Precedence

```
expression -> ternary 
ternary    -> equality
              | equality "?" expression ":" ternary ;

equality   -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term       -> factor ( ( "-" | "+" ) factor )* ;
factor     -> unary ( ( "/" | "*" ) unary )* ;

unary      -> ( "!" | "-" ) unary
              | primary ;

primary    -> NUMBER | STRING | "true" | "false" | "nil"
              | "(" expression ")" ;
```
