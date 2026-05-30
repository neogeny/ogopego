"""ogopego — a PEG parser generator for Python, powered by Go.

Provides compile(), parse(), and parse_file() with a TatSu-compatible
API that compiles PEG grammars and parses input text, delegating all
work to a bundled Go binary.

Usage:

    import ogopego

    grammar = '''
        start = number $ ;
        number = digit+ ;
        digit = /[0-9]/ ;
    '''

    # compile a grammar (cached by content hash)
    compiled = ogopego.compile(grammar)

    # parse input text
    tree = ogopego.parse(grammar, "42")
    # tree == {"number": {"digit": ["4", "2"]}}

    # parse a file
    tree = ogopego.parse_file(grammar, "/path/to/input.txt")

After installation the bundled Go binary is also available as the ``ogo``
command-line tool:

    $ ogo run grammar.ebnf input.txt --json
"""

from ogopego._ogopego import OgoError, compile, parse, parse_file

__all__ = ["OgoError", "compile", "parse", "parse_file"]
