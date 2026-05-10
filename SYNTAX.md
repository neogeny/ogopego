orphan

:

# Grammar Syntax

**铁修 TieXiu** grammars use an extension of
the classic
[EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form) syntax.
The classic variations of
[EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form)
(Tomassetti, EasyExtend, Wirth) and [ISO
EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form) are also
supported as input grammar format.

## Rules

A grammar consists of a sequence of one or more rules of the form:

```ebnf
rulename: <expre>
```

If a *name* collides with a [Python](http://python.org) keyword or
builtin, an underscore (`_`) will be appended to it on the generated
parser.

Rule names that start with an uppercase character:

```ebnf
FRAGMENT: /[a-z]+/
```

*do not* advance over whitespace or comments before beginning to parse. This feature
becomes handy when defining complex lexical elements, as it allows
breaking them into more than one rule.

The parser returns an
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) value for each
rule depending on what was parsed:

- A single value
- A list of [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree)
- A mapping for rules with named elements
- An node with metadata

This is the `enum Tree` that defines the posible results of a parse:

```rust
pub enum Tree {
    Void,                          // The result of the () expression
    Text(Box<str>),                // Tokens or patterns
    List(Box<[Tree]>),             // Sequences of expressions
    Map(Box<TreeMap>),             // A mapping of named elements

    Node {                         // The result of parsing a rule
        meta: NodeMetaRef,         // Metadata
        tree: Box<Tree>            // The result of parsing a rule
    },
    // ...
```


## Expressions

The expressions, in reverse order of operator precedence, can be any of
the following.

> **note**
>
> Because **TieXiu** now supports free-form [EBNF][], there must not be empty lines within expressions.

[EBNF]: https://en.wikipedia.org/wiki/Extended_Backus–Naur_form


### `# comment ...`

[Python](http://python.org)-style end-of-line comments are allowed.

### `// comment ...`

[Java](https://en.wikipedia.org/wiki/Java)-style end-of-line comments
are allowed.

### `/* ... */`

[EBNF](https://en.wikipedia.org/wiki/Extended_Backus–Naur_form)-style
multi-line comments are allowed.

### `e1 | e2`

Choice. Match either `e1` or `e2`.

A `|` may be used before the first option if desired:

```ebnf
choices:
  | e1
  | e2
  | e3
```

### `e1 e2`

Sequence. Match `e1` and then match `e2`.

### `( e )`

Grouping. Match `e`. For example: `('a' | 'b')`.

### `(?: e )`

A non-capturing group. Like in a `()` group, match `e`, but this time do
not capture what was parsed. For example:

``` {.apl force=""}
start: header (?: delimiter ) body

header: /[A-Z]+/

delimiter: /[:,-]+/

body: /[a-z]+/
```

### `[ e ]`

Optionally match `e`.

### `{ e }` or `{ e }*`

Closure. Match `e` zero or more times. The
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) returned for a
closure is always a `list`.

### `{ e }+`

Positive closure. Match `e` one or more times. The
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) is always a
`list`.

### `{}`

Empty closure. Match nothing and produce an empty `list` as
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree).

### `~`

The *cut* expression: commit to the current option and prevent other
options from being considered even when what follows fails to parse.

In this example, other options won\'t be considered if a parenthesis is
parsed:

```ebnf
atom:
    | '(' ~ @:expre ')'
    | int
    | bool
```

Cut expression may be used anywhere. The effect of `~` is scoped to the
nearest enclosing brackets (*group*, *optional*, *closure*), the
enclosing *choice*, or the enclosing *rule*.

On the scoping of *cut*, consider these theoretical equivalences about
implicit choices in some expressions:

``` {.apl force=""}
A →  α   ≡ A → α | ⊥
A → [x]  ≡ A → B, B → x | ε
A → {x}  ≡ A → B, B → xB | ε
A → {x}+ ≡ A → B, B → xB | x
```

This is a common use of `~`. The *closure* doesn\'t allow a partial
assignment expressions to escape it:

```ebnf
parameters: ','.{name '=' ~ expression}
```

### `s%{ e }+`

Positive join. Inspired by [Python](http://python.org)\'s `str.join()`,
it parses the same as this expression:

```ebnf
e {s ~ e}
```

yet the result is a single list of the form:

```ebnf
[e, s, e, s, e, ...]
```

Use grouping if [s]{.title-ref} is more complex than a *token* or a
*pattern*:

```ebnf
(s t)%{ e }+
```

### `s%{ e }` or `s%{ e }*`

Join. Parses the list of `s`-separated expressions, or the empty
closure. It\'s equivalent to:

```ebnf
s%{e}+|{}
```

### `s.{ e }+`

Positive *gather*. Like *positive join*, but the separator is not
included in the resulting
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree).

### `s.{ e }` or `s.{ e }*`

*Gather*. Like the *join*, but the separator is not included in the
resulting [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree).
It\'s equivalent to:

```ebnf
s.{e}+|{}
```

### `&e`

Positive lookahead. Succeed if `e` can be parsed, but do not consume any
input.

### `!e`

Negative lookahead. Fail if `e` can be parsed, and do not consume any
input.

### `'text'` or `"text"`

Match the token *text* within the quotation marks.

Note that if *text* is alphanumeric, then **TieXiu** will check that the character following the token is not alphanumeric. This is done to prevent tokens like *IN* matching when the text ahead is *INITIALIZE*. This feature can be turned off by passing `nameguard=False` to the `Parser` or the `Buffer`, or by using a pattern expression (see below) instead of a token expression. The `@@nameguard` and `@@namechars` directives may be specified in the grammar for the same effect:

```ebnf
@@nameguard :: False
```

or to specify additional characters that should also be considered part
of names:

```ebnf
@@namechars :: '$-.'
```

### `r'text'` or `r"text"`

Match the token *text* within the triple quotation marks, interpreting
*text* like [Python](http://python.org)\'s [raw string
literal](https://docs.python.org/3/reference/lexical_analysis.html#string-and-bytes-literals)s.
This is useful

### `'''text'''` or `"""text"""`

A multi-line version of `token`. The text within the triple quotation
marks is stripped of trailing whitespace, and the common indentation is
removed as to help with pretty formatting grammars without having to
worry about indentation. The resulting text is matched as a token.

### `/regexp/`

Also `?"regexp"` or `?'regexp'`, The *pattern* expression. Match the
[Python](http://python.org) regular expression `regexp` at the current
text position. Unlike other expressions, this one does not advance over
whitespace or comments. For that, place the `regexp` as the only term in
its own rule.

The *regex* is interpreted as a [Python](http://python.org) [raw string
literal](https://docs.python.org/3/reference/lexical_analysis.html#string-and-bytes-literals)
and passed the [Python](http://python.org)
[re](https://docs.python.org/3.4/library/re.html) module using `match()`
at the current position in the text. The returned
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) has the
semantics of `re.findall(pattern, text)[0]` (a [tuple]{.title-ref} if
there is more than one group), so use `(?:)` for groups that should not
be in the resulting
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree).

Consecutive *patterns* are concatenated to form a single one.

### `/./`

The *any* expression, matches the next position in the input. It works
exactly like the `?'.'` pattern, but is implemented at the lexical
level, without regular expressions.

### `->e`

The \"\*skip to\*\" expression, useful for writing *recovery* rules.

The parser will advance over input, one character at time, until `e`
matches. Whitespace and comments will be skipped at each step. Advancing
over input is done efficiently, with no regular expressions involved.

The expression is equivalent to:

```ebnf
{ !e /./ } e
```

A common form of the expression is `->&e`, which is equivalent to:

```ebnf
{ !e /./ } &e
```

This is an example of the use of the \"\*skip to\*\" expression for
recovery:

```ebnf
statement:
    | if_statement
    # ...

if_statement:
    | 'if' condition 'then' statement ['else' statement]
    | 'if' statement_recovery

statement_recovery: ->&statement
```

### `` `constant ``\`

Match nothing, but behave as if `constant` had been parsed.

Constants can be used to inject elements into the concrete and abstract
syntax trees, perhaps avoiding having to write a semantic action. For
example:

```ebnf
boolean_option: name ['=' (boolean|`true`) ]
```

If the text evaluates to a Python literal (with `ast.literal_eval()`),
that will be the returned value. Otherwise, string interpolation in the
style of `str.format()` over the names in the current
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) is applied for
*constant* elements. Occurrences of the `{` character must be escaped to
`\{` if they are not intended for interpolation. A *constant* expression
that hast type `str` is evaluated using:

``` {.python force=""}
eval(f'{"f" + repr(text)}', {}, ast)
```

### ``[constant]{.title-ref}``

A multi-line version of `` `constant ``\`.

### `` ^`constant ``\`

Also `^`[constant]{.title-ref}``. An alert. There will be no token
returned by the parser, but an alert will be registed in the parse
context and added to the current node\'s `parseinfo`.

The `^` character may appear more than once to encode the *alert level*:

```ebnf
assignment: identifier '=' (
    | value
    | -> &';' ^^^`could not parse value in assignment to {identifier}`
)
```

### `rulename`

Invoke the rule named `rulename`. To help with lexical aspects of
grammars, rules with names that begin with an uppercase letter will not
advance the input over whitespace or comments.

### `>rulename`

The include operator. Include the *right hand side* of rule `rulename`
at this point.

The following set of declarations:

```ebnf
includable: exp1

expanded: exp0 >includable exp2
```

Has the same effect as defining *expanded* as:

```ebnf
expanded: exp0 exp1 exp2
```

Note that the included rule must be defined before the rule that includes it.

### `()`

The empty expression. Succeed without advancing over input. Its value is
the empty tuple `()`.

### `!()`

The *fail* expression. This is `!` applied to `()`, which always fails,
and thus has no value. This is useful to prevent the parser from
accepting certain input, for example in placeholders for rules that are
not yet defined.

### `name=e` or `name:e`

Add the result of `e` to the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) using `name` as
key. If `name` collides with any attribute or method of `dict`, or is a
[Python](http://python.org) keyword, an underscore (`_`) will be
appended to the name.

> **Note**
>
> `name` is bound in the *option* in which it appears, or in the rule when there are no options. When options define different names, only the names in the option that parses will be present in the resulting [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree). A `name` will be bound to `None` when the expression `e` fails to parse. For `name` used in enclosing expressions like *group*, *optional*, or *closure*, `name` will be bound in the rule-level [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) only if the complete enclosure parses (they have local scope, and are transferred to the outer scope only on success). The same criteria applies to expressions nested to any level.

When there are no named items in a rule or choice, the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) consists of the
elements parsed by the rule, either a single item or a `list`. This
default behavior makes it easier to write simple rules:

```ebnf
number: /[0-9]+/
```

Without having to write:

```ebnf
number: number=/[0-9]+/
```

When a rule has named elements, the unnamed ones are excluded from the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) (they are
ignored).

### `name+=e` or `name+:e`

Add the result of `e` to the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) using `name` as
key. Force the AST entry to be a `list` even if only one element is
added. Collisions with `dict` attributes or [Python](http://python.org)
keywords are resolved by appending an underscore to `name`.

### `=e` or `@:e`

The override operator. Make the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) for the
complete rule or choice be the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) for `e`.

> **note**
>
>As with `name=e`, the effect of `=e` is scoped to the enclosing *option*, *group*, *optional*, or *closure*, and will apply only when the enclosure parses successfully.

This expression is useful to recover only part of the right hand side of
a rule without the need to name it, or add a semantic action. This is a
typical use of the override operator:

```ebnf
subexp: '(' =expre ')'
```

The [AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) returned
for the `subexp` rule will be the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) recovered from
invoking `expre`.

### `+=e` or `@+:e`

Like `=e`, but make the
[AST](http://en.wikipedia.org/wiki/Abstract_syntax_tree) always be a
`list`.

This operator is convenient in cases such as this, in which the
delimiting tokens are of no interest.

```ebnf
arglist: '(' +=arg {',' +=arg}* ')'
```

### `$`

The *end of text* symbol. Verify that the end of the input text has been
reached.

### `$->`

The *end of line* symbol. Verify that the end of the current line has
been reached. This is useful for parsing line-based formats, such as
configuration files, or for parsing comments.

The `$->` (EOL) expression will consume the whitespace up to and including the next line break, using the Python semantics of `os.linesep`. The match interprets whitespace using the Python definition as implemented by `str .isspace()`, so beware when a particular definition of *whitespace* is part of the language being parsed.

Comments, as defined for the grammar, will also be skipped by the `$->` expression in search of a newline, which means that newlines consummed by the comments patterns will not be *"seen"* by `$->`.


## Rules with Arguments

Rules may specify a list of arguments and keyword arguments by enclosing
them in square brackets right after the rule name:

```ebnf
addition[Add, op='+']: addend '+' addend
```

The arguments values are fixed at grammar-compilation time.

## Based Rules

Rules may extend rules defined before by using the `<` operator.

The following set of declarations:

```ebnf
base[Param]: exp1

extended < base: exp2
```

Has the same effect as defining *extended* as:

```ebnf
extended[Param]: exp1 exp2
```

Parameters from the *base rule* are copied to the new rule if the new
rule doesn\'t define its own. Repeated inheritance should be possible,
but it *hasn\'t been tested*.

## Memoization

**TieXiu** generates *packrat* parsers.
The result of parsing a rule at a given position in the input is
*memoized*. The next time the parser visits the same input position,
with the same rule, the memoized result is returned and the input
advanced accordingly.

Some rules should not be memoized. For example, rules that may succeed
or not depending on the associated semantic action *should not* be
memoized.

The `@nomemo` decorator turns off memoization for a particular rule:

```ebnf
@nomemo
INDENT: ()

@nomemo
DEDENT: ()
```

## Rule Overrides

A grammar rule may be redefined by using the `@override` decorator:

```ebnf
start: ab $

ab: 'xyz'

@override
ab: @:'a' {@:'b'}
```

## Grammar Name

The name to be used in any generated artifacts can be specified within the grammar using the `@@grammar` directive:

```ebnf
@@grammar :: MyLanguage
```

## Whitespace

By default, **TieXiu** generated parsers skip the usual whitespace characters with the regular expression `r'\s+'`, but you can change that behavior.

Whitespace may be specified within the grammar using the `@@whitespace`
directive in the grammar:

```ebnf
@@whitespace :: /[\t ]+/
```

or:

```ebnf
@@whitespace :: None
```

If no `whitespace` or `@@whitespace` is specified, **TieXiu** will use `r'(?m)\s+'` as a default. Use `None` to have *no whitespace definition*.

The character string is converted into a regular expression character
set before starting to parse.

If you do not define any whitespace characters, then you will have to
handle whitespace in your grammar rules (as it\'s often done in
[PEG][] parsers):

[PEG]: http://en.wikipedia.org/wiki/Parsing_expression_grammar

## Case Sensitivity

You may specify case insensitivity within the grammar using the `@@ignorecase` directive:

```ebnf
@@ignorecase :: True
```

The change will affect token matching, but not pattern matching. Use
`(?i)` in patterns that should ignore case.

## Comments

Patterns for comments may be specified within a grammar using the `@@comments` and `@@eol_comments` directives:

```ebnf
@@comments :: /\(\*.*?\*\)/
@@eol_comments :: /#.*?$/
```

## Reserved Words and Keywords

Some languages must reserve the use of certain tokens as valid
identifiers because the tokens are used to mark particular constructs in
the language. Those reserved tokens are known as [Reserved
Words](https://en.wikipedia.org/wiki/Reserved_word) or
[Keywords](https://en.wikipedia.org/wiki/Reserved_word)

**TieXiu** provides support for preventing the use of [keywords](https://en.wikipedia.org/wiki/Reserved_word) as identifiers
though the `@@keyword` directive,and the `@name` decorator.

A grammar may specify reserved tokens providing a list of them in one or
more `@@keyword` directives:

```ebnf
@@keyword :: if endif
@@keyword :: else elseif
```

The `@name` decorator checks that the result of a grammar rule does not match a token defined as a [keyword](https://en.wikipedia.org/wiki/Reserved_word):

```ebnf
@name
identifier: /(?!\d)\w+/
```

Note that the rule decorated with `@name` must produce a single string
as result (no named expressions that will produce a dict, and no rule
arguments).

In some situations a token is reserved only in a specific context. In
those cases, a negative lookahead will prevent the use of that token:

```ebnf
statements: {!'END' statement}+
```

## Left Recursion

**TieXiu** supports left recursion in [PEG][] grammars. The algorithm used is [Warth et al][]'s. Left recursion support is enabled by default.

[PEG]: http://en.wikipedia.org/wiki/Parsing_expression_grammar
[Warth et al]: http://www.vpri.org/pdf/tr2007002_packrat.pdf

Left recursion can be turned *on* or *off* from within the grammar using
the `@@left_recursion` directive:

```ebnf
@@left_recursion :: False
``


## Grammar Syntax

```console
start ●─grammar─■

grammar ●─ title=(`TATSU`) ──┬→─────────────────────────────────┬── rules+=(rule) ──┬→─────────────────────────────┬──⇥＄
                             ├→──┬─ directives+=(directive) ─┬──┤                   ├→──┬─ rules+=(rule) ───────┬──┤
                             │   └─ keywords+=(keyword) ─────┘  │                   │   └─ keywords+=(keyword) ─┘  │
                             └─────────────────────────────────<┘                   └─────────────────────────────<┘

directive ●─"@@"─ !["keyword"] ✂ ──┬─ name=(──┬─"comments"─────┬─)  ✂ "::" ✂  value=(regex) ────────────┬─ ✂ ─■
                                   │          └─"eol_comments"─┘                                        │
                                   ├─ name=("whitespace")  ✂ "::" ✂  value=(──┬─regex───┬─) ────────────┤
                                   │                                          ├─string──┤               │
                                   │                                          ├─"None"──┤               │
                                   │                                          ├─"False"─┤               │
                                   │                                          └─`None`──┘               │
                                   ├─ name=(──┬─"nameguard"──────┬─)  ✂ ──┬─"::" ✂  value=(boolean) ─┬──┤
                                   │          ├─"ignorecase"─────┤        └─ value=(``) ─────────────┘  │
                                   │          ├─"left_recursion"─┤                                      │
                                   │          ├─"parseinfo"──────┤                                      │
                                   │          └─"memoization"────┘                                      │
                                   ├─ name=("grammar")  ✂ "::" ✂  value=(word) ─────────────────────────┤
                                   └─ name=("namechars")  ✂ "::" ✂  value=(string) ─────────────────────┘

keywords ●───┬─keyword─┬───■
             └────────<┘

keyword ●─"@@keyword" ✂ "::"──┬─ +=(──┬─word───┬─) ─ ![──┬─":"─┬─]─┬───■
                              │       └─string─┘         └─"="─┘   │
                              └───────────────────────────────────<┘

params ●─ +=(first_param) ──┬→────────────────────────────┬───■
                            ├→"," +=(literal) ─ !["="] ✂ ─┤
                            └────────────────────────────<┘

first_param ●───┬─path────┬──■
                └─literal─┘

kwparams ●───┬─"," │ pair─┬───■
             └───────────<┘

the_params_at_last ●───┬─ kwparams=(kwparams) ────────────────────────┬──■
                       ├─ params=(params) "," ✂  kwparams=(kwparams) ─┤
                       └─ params=(params) ────────────────────────────┘

paramdef ●───┬─"[" ✂  >(the_params_at_last) "]"─┬──■
             ├─"(" ✂  >(the_params_at_last) ")"─┤
             └─"::" ✂  params=(params) ─────────┘

rule ●─ decorators=(──┬→──────────┬──)  name=(name)  ✂ ──┬─→ >(paramdef) ─┬───┬─→"<" ✂  base=(known_name) ─┬─ ∅ /=|::=|:=?/─ ✂  exp=(expre) ENDRULE ✂ ─■
                      ├→decorator─┤                      └─→──────────────┘   └─→──────────────────────────┘
                      └──────────<┘

ENDRULE ●───┬─DEDENT─┬─
            ├─BLANK──┤
            ├─";"────┤ ─┘
            └─⇥＄

DEDENT ●─EOL/\S/──■

BLANK ●─EOLEOL─■

EOL ●─/(?m)[ \t]*$/─/(?m)(?:\r?\n|\r)?/──■

decorator ●─"@"─ !["@"] ✂  =(──┬─"override"─┬─) ─■
                               ├─"name"─────┤
                               ├─"isname"───┤
                               └─"nomemo"───┘

pair ●─ +=(word) "=" ✂  +=(literal) ─■

expre ●───┬─choice───┬──■
          └─sequence─┘

choice ●───┬─→"|" ✂ ─┬─ +=(option) ──┬─"|" ✂  +=(option) ─┬───■
           └─→───────┘               └───────────────────<┘

option ●─sequence─■

sequence ●───┬── &[element","]──┬─"," │ element─┬───┬──■
             │                  └──────────────<┘   │
             └───┬── ![ENDRULE]element─┬────────────┘
                 └────────────────────<┘

element ●───┬─named────────┬──■
            ├─term─────────┤
            ├─override─────┤
            └─rule_include─┘

rule_include ●─">" ✂  =(known_name) ─■

named ●───┬─named_list───┬──■
          └─named_single─┘

named_list ●─ name=(name) /\+[:=]/─ ✂  exp=(term) ─■

named_single ●─ name=(name) /[:=]/─ ✂  exp=(term) ─■

override ●───┬─override_list──────────────┬──■
             ├─override_single────────────┤
             └─override_single_deprecated─┘

override_list ●─/\+=|@\+:/─ ✂  =(term) ─■

override_single ●─/=|@:/─ ✂  =(term) ─■

override_single_deprecated ●─"@" ✂  =(term) ─■

term ●───┬─gather─────────────┬──■
         ├─join───────────────┤
         ├─left_join──────────┤
         ├─right_join─────────┤
         ├─empty_closure──────┤
         ├─positive_closure───┤
         ├─closure────────────┤
         ├─optional───────────┤
         ├─atom───────────────┤
         ├─void───────────────┤
         ├─skip_to────────────┤
         ├─lookahead──────────┤
         ├─negative_lookahead─┤
         ├─cut────────────────┤
         └─cut_deprecated─────┘

group ●── !["(?:"]"(" ✂  =(expre) ")" ✂ ─■

skip ●─"(?:" ✂  =(expre) ")" ✂ ─■

gather ●── &[atom".{"] ✂ ──┬─positive_gather─┬──■
                           └─normal_gather───┘

positive_gather ●─ sep=(atom) ".{" exp=(expre) "}"/(?!\+=)[+-]/─ ✂ ─■

normal_gather ●─ sep=(atom) ".{" ✂  exp=(expre) "}"──┬─→"*" ✂ ─┬─ ✂ ─■
                                                     └─→───────┘

join ●── &[atom"%{"] ✂ ──┬─positive_join─┬──■
                         └─normal_join───┘

positive_join ●─ sep=(atom) "%{" exp=(expre) "}"/(?!\+=)[+-]/─ ✂ ─■

normal_join ●─ sep=(atom) "%{" ✂  exp=(expre) "}"──┬─→"*" ✂ ─┬─ ✂ ─■
                                                   └─→───────┘

left_join ●─ sep=(atom) "<{" ✂  exp=(expre) "}"/(?!\+=)[+-]/─ ✂ ─■

right_join ●─ sep=(atom) ">{" ✂  exp=(expre) "}"/(?!\+=)[+-]/─ ✂ ─■

positive_closure ●───┬─"{" =(expre) "}"/(?!\+=)[+-]/─ ✂ ─┬──■
                     └─ =(atom) /(?!\+=)[+]/─ ✂ ─────────┘

closure ●───┬─"{" =(expre) "}"──┬─→"*"─┬─ ✂ ─┬──■
            │                   └─→────┘     │
            └─ =(atom) "*" ✂ ────────────────┘

empty_closure ●─"{}" ✂  =( ∅ ) ─■

optional ●───┬─"[" ✂  =(expre) "]" ✂ ───────────┬──■
             └─ =(atom) ─ ![──┬─"?\""─┬─]"?" ✂ ─┘
                              ├─"?'"──┤
                              └─"?/"──┘

lookahead ●─"&" ✂  =(term) ─■

negative_lookahead ●─"!" ✂  =(term) ─■

skip_to ●─"->" ✂  =(term) ─■

atom ●───┬─token────┬──■
         ├─call─────┤
         ├─dot──────┤
         ├─pattern──┤
         ├─group────┤
         ├─eol──────┤
         ├─eof──────┤
         ├─skip─────┤
         ├─alert────┤
         └─constant─┘

call ●─word─■

void ●─"()" ✂ ─■

fail ●─"!()" ✂ ─■

cut ●─"~" ✂ ─■

cut_deprecated ●─">>" ✂ ─■

known_name ●─name ✂ ─■

name ●─word─■

constant ●── &["`"]──┬─/(?ms)```((?:.|\n)*?)```/──┬──■
                     ├─"`" =(literal) "`"─────────┤
                     └─/`(.*?)`/──────────────────┘

alert ●─ level=(/\^+/─)  message=(constant) ─■

token ●───┬─string─────┬──■
          └─raw_string─┘

literal ●───┬─string─────┬──■
            ├─raw_string─┤
            ├─boolean────┤
            ├─word───────┤
            ├─hex────────┤
            ├─float──────┤
            ├─int────────┤
            └─null───────┘

string ●── &[──┬─"\""─┬─]──┬─multiline_string─┬──■
               └─"'"──┘    ├─singlequoted─────┤
                           └─doublequoted─────┘

singlequoted ●─SINGLEQUOTED─■

doublequoted ●─DOUBLEQUOTED─■

raw_string ●─/r(?=["'])/─ =(STRING) ─■

STRING ●───┬─SINGLEQUOTED─┬──■
           └─DOUBLEQUOTED─┘

SINGLEQUOTED ●─/'((?:[^'\n]|\\'|\\\\)*?)/─ ✂ ─■

DOUBLEQUOTED ●─/"((?:[^"\n]|\\"|\\\\)*?)"/─ ✂ ─■

multiline_string ●───┬─/(?ms)'''((?:\\\\|\\.|(?!''').)*?)/─ ✂ ────┬──■
                     └─/(?ms)"""((?:\\\\|\\.|(?!""").)*?)"""/─ ✂ ─┘

hex ●─/0[xX](?:\d|[a-fA-F])+/──■

float ●─/[-+]?(?:\d+\.\d*|\d*\.\d+)(?:[Ee][-+]?\d+)?/──■

int ●─/[-+]?\d+/──■

path ●─/(?!\d)\w+(?:::(?!\d)\w+)+/──■

word ●─/(?!\d)\w+/──■

dot ●─"/./"─■

pattern ●─regex─■

regex ●───┬─deprecated_regex───────────────┬──■
          └── !["?/"]──┬─REGEX──────────┬──┘
                       └─"?" =(STRING) ─┘

REGEX ●── &["/"]/(?ms)/((?:[^/\\]|\\/|\\.)*)//─ ✂ ─■

deprecated_regex ●─"?/" ✂  =(/(?ms)((?:[^/\\]|\\/|\\.)*)/─)  ✂ "/?"─■

boolean ●───┬─"True"──┬──■
            └─"False"─┘

null ●─"None"─■

eof ●─"$" ✂ ─■

eol ●─"$->"─■


```
