# goitr

goitr is a simple interpreter for educational purposes, mainly
referencing the "Crafting Interpreter" book but with my own implementation
to port into go.

The project's purpose is for me to learn more about go and how programming
languages work in general, so bugs and imperfect code all to be
expected. It's not a full-fledge product anyway.

The remaining sections are mainly my notes and understandings on the
topic of interpreters.

## Process

<source_code> -> lexing -> <tokens> -> parsing -> <AST> -> static
analysis -> <intermediate_representation> -> optimize

### Asides (read more)

IR (intermediate representation):
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

## License

See [LICENSE](./LICENSE) for more information.
