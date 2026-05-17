# ⻰OgoPEGo

A PEG parser generator in Go.

**⻰OgoPEGo** is the Go sibling of [竜TatSu] (Python) and [铁修TieXiu] (Rust).
It is functionally complete and passes the same test suite as its siblings.

[竜TatSu]: https://tatsu.readthedocs.io/
[铁修TieXiu]: https://github.com/neogeny/TieXiu

## Documentation

Refer to the [竜TatSu documentation] for grammar syntax, semantics, and usage.
The local [SYNTAX.md](SYNTAX.md) describes the grammar format.

[竜TatSu documentation]: https://tatsu.readthedocs.io/

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

## Non-Features

Most features of **竜TatSu** are available in **⻰OgoPEGo**. Some features have not yet been implemented, and a few never will:

* [ ] Generation of synthetic classes from grammar parameters will not be implemented in Go.
* [ ] Generation of source code with an object model for deifinitions in the grammar may be implemented if a way is found to make the parser or postprocessing bind the Tree output of a parse to the model.
* [ ] Code generation of a parser recently moved in **竜TatSu** to the loading of a model of the Grammar and using it as parser. Although the generated procedural parser may produce 1.3x increased throughput in Python, supporting generated code is hard and it complicates the internal interfaces. For Go, **⻰OgoPEGo** _alreay knows_ how to load _fast_ a Grammar model from **竜TatSu** JSON. A generated copy of the grammar model constructor could be precompiled by Go.
* [ ] Parsing of boolean and numeric values happens in **竜TatSu** through synthetic actions, which call the constructors for those types passing the parsed strings. For **⻰OgoPEGo** the preferred way of transformig a tree (semantics) is through post-processing (folding), but basic numeric types and booleans could be supported.
* [ ] Semantic actions (transformations) during parse are not implemented. Python is friendly to objects of type `Any`, so semantic actions during parse in **竜TatSu** can produce a _tree_ of any type. Go is different, and trying to produce structures of type `any` is not idiomatic. The result of a parse is a well-defined Tree which is a small-enough interface that writing a walker for it is easy, so type transformations can be done in postprocessing by folding.
* [ ] Interpolation and evaluation of _\`constant\`_ expressions hasn't had any known use cases with **竜TatSu**. They will not be implemented in **⻰OgoPEGo** until a use case appears.


## License

Licensed under either of:

* Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

### Contribution

Unless explicitly stated otherwise, any contribution intentionally submitted
for inclusion in the work, as defined in the Apache-2.0 license, shall be
dual-licensed as above, without any additional terms or conditions.
