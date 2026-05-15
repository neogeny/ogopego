# ⻰OgoPEGo

A PEG parser generator in Go.

**⻰OgoPEGo** is the Go sibling of [TatSu] (Python) and [TieXiu] (Rust).
It is functionally complete and passes the same test suite as its siblings.

[TatSu]: https://tatsu.readthedocs.io/
[TieXiu]: https://github.com/neogeny/TieXiu

## Documentation

Refer to the [TatSu documentation] for grammar syntax, semantics, and usage.
The local [SYNTAX.md](SYNTAX.md) describes the grammar format.

[TatSu documentation]: https://tatsu.readthedocs.io/

## Installation

```bash
go install github.com/neogeny/ogopego/cmd/ogo@latest
```

## Library API

```go
import "github.com/neogeny/ogopego/api"

// Compile a grammar string into a Grammar object.
g, err := api.Compile(grammar, cfg)

// Parse input with a compiled Grammar.
tree, err := api.ParseInput(g, input, cfg)

// Compile and parse in one step.
tree, err := api.ParseGrammar(grammar, cfg)

// Compile to JSON-compatible output.
json, err := api.CompileToJSON(grammar, cfg)

// Parse input to JSON-compatible output.
json, err := api.ParseInputToJSON(g, input, cfg)

// JSON roundtrip via peg package.
jsonStr := peg.SerializeGrammar(g)
g2, err := peg.ParseGrammar([]byte(jsonStr))
```

### Grammar object

```go
import "github.com/neogeny/ogopego/peg"

// A compiled grammar. Create one with api.Compile.
type Grammar struct {
    Name       string            // grammar name
    Directives *asjson.OrderedMap // @@directives
    Keywords   []string           // @@keyword declarations
    Rules      []*Rule            // grammar rules
    Analyzed   bool               // true after Initialize()
}

// Parse input text with this grammar (use api.ParseInput).
result, err := api.ParseInput(g, text, cfg)

// Prepare grammar for parsing (link rules, detect left recursion).
err := g.Initialize()

// Serialize.
jsonStr := g.AsJSONStr()             // indented JSON
jsonStr := peg.SerializeGrammar(g)   // clean JSON (recommended)
data, err := peg.ParseGrammar([]byte(jsonStr))  // deserialize

// Display.
fmt.Println(g.PrettyPrint())    // EBNF pretty-print
fmt.Println(g.Railroads())      // railroad diagram
```

## CLI

The `ogo` CLI is a convenience for testing grammars and examining
output formats:

```bash
ogo run grammar.json input.txt        # parse input files
ogo boot --pretty                     # inspect boot grammar
ogo grammar grammar.ebnf --railroads  # diagram grammar
```

Use `ogo --help` for details.

## License

Licensed under either of:

* Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

### Contribution

Unless explicitly stated otherwise, any contribution intentionally submitted
for inclusion in the work, as defined in the Apache-2.0 license, shall be
dual-licensed as above, without any additional terms or conditions.
