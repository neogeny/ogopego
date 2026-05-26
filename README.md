# ⻰OGoPEGo

A PEG parser generator in Go.

**⻰OGoPEGo** is the Go sibling of [竜TatSu] (Python) and [铁修TieXiu] (Rust).
It is functionally complete and passes the same test suite as its siblings.

[竜TatSu]: https://tatsu.readthedocs.io/
[铁修TieXiu]: https://github.com/neogeny/TieXiu

Refer to the [竜TatSu documentation] for grammar syntax, semantics, and usage. The local [SYNTAX.md](SYNTAX.md) describes the grammar format.

[竜TatSu documentation]: https://tatsu.readthedocs.io/

The CLI tool is a great way to explored the features offered by the library:

```bash
$ ogo --help
Usage: ogo <command> [flags]

ogopego: A PEG parser generator in Go

Flags:
  -h, --help             Show context-sensitive help.
  -o, --output=STRING    Output to a file instead of stdout
  -C, --color="auto"     Control colorized output for API results
  -t, --trace            Display a detailed trace of the parsing process
  -v, --version          Print version information

Commands:
  run <grammar> <inputs> ... [flags]
    Execute a grammar against one or more input files

  boot [flags]
    The internal boot grammar

  grammar <grammar> [flags]
    Grammar transformations

Run "ogo <command> --help" for more information on a command.
```

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

## Features

* [x] Generation of source code with an object model for deifinitions in the grammar is complete. The model generator defines the types specified in the input grammar and creates the transformation from the `Tree` result of a call to `Parse()` to the object model.
* [x] Code generation of a parser recently moved in **竜TatSu** to the loading of a model of the Grammar and using it as parser. **⻰OGoPEGo** is cabable of generating a model of the `Grammar` constructor for a grammar that can be compiled by Go for blazing boot times.
* **⻰OGoPEGo** also knows how to load _fast_ a `Grammar` model from **竜TatSu**-format JSON.
* [x] Semantic actions (transformations) during parse are implemented through a `SemanticsFunc(Tree) Tree` configuration entry. Transformations are limited to `Tree->Tree`, so a walker must be used as post-processor if a different AST type is desidred. The AST model generation available through the CLI tool generates such a walker. 
* [ ] Interpolation and evaluation of _\`constant\`_ expressions hasn't had any known use cases with **竜TatSu**. They will not be implemented in **⻰OGoPEGo** until a use case appears.


## License

Licensed under the Apache License, Version 2.0 ([LICENSE](LICENSE) or http://www.apache.org/licenses/LICENSE-2.0).

### Contribution

Unless explicitly stated otherwise, any contribution intentionally submitted
for inclusion in the work, as defined in the Apache-2.0 license, shall be
licensed as above, without any additional terms or conditions.
