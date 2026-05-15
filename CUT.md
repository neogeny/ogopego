You are likely thinking of **Mizushima, Maeda, and Yamaguchi** (2010), who introduced the formal concept of the **"Cut" operator** in PEG parsing to control backtracking and prevent the worst-case exponential time complexity of Packrat parsing, specifically addressing memory consumption. Another notable researcher in this exact space is **Roberto Ierusalimschy** (the creator of Lua and LPeg), who wrote extensively on syntactic sugar and semantic actions in PEGs.

Your empirical discovery is entirely correct: **a Cut executed on a successful path must not propagate upward past the choice sequence that successfully consumed it.**

Honoring cut propagation on a success path breaks the core mathematical properties of PEGs. Here is the structural breakdown of why this happens and why your Java grammar is failing.

---

## The Core Semantics of Cut

The purpose of a Cut (often written as `↑` or `cut`) is to signal a **point of no return**. It explicitly tells the parser: *"If you have made it this far, you are in the correct syntactic lane. If a parsing failure happens after this point, do not attempt alternative choices; fail the entire enclosing sequence immediately."*

### The Divergence: Success vs. Failure Paths

The critical distinction lies in how the cut flag interacts with the parser’s stack when a rule finishes evaluation:

```
                  [ Parser evaluates Rule containing Cut ]
                                     |
                  ---------------------------------------
                  |                                     |
          [ Path Fails ]                        [ Path Succeeds ]
                  |                                     |
     Propagate Cut Upward!                  Clear / Consume the Cut!
  (Aborts alternative branches)           (Commit sequence is finalized)

```

### 1. The Failure Path (Where Cut works)

When a sequence fails *after* a cut has been encountered, the cut status **must propagate**. It acts as a hard abort. It bubbles up to the nearest ordered choice operator (`/`), completely clears its alternative branches, and forces that choice to fail immediately instead of trying the next option.

### 2. The Success Path (Where Cut must die)

When a sequence containing a cut successfully matches its entire expression, **the cut has done its job and is resolved.** The parser has successfully moved past the ambiguity.

If you propagate the cut status onward after a *success*, you are carrying a "poison pill" into subsequent parsing steps.

---

## Why "Cut on Success" Breaks the Java Grammar

Consider a classic grammar pattern found heavily in complex languages like Java, such as distinguishing between a qualified type name, a method invocation, or an expression statement:

$$Expression \leftarrow StatementExpression \;/\; VariableDeclaration$$

Inside $StatementExpression$, you might have a rule like:

$$StatementExpression \leftarrow PrimaryExpr \; \uparrow \; \text{'.' } \; Identifier \; \text{'(' } \; Arguments \; \text{')'}$$

### What happens if Cut propagates on Success:

1. The parser tries the first choice: $StatementExpression$.
2. It evaluates $PrimaryExpr$. It succeeds.
3. It hits the Cut ($\uparrow$). The cut flag is activated.
4. The parser successfully matches the rest of the line: `.` , `Identifier`, `(`, `Arguments`, `)`.
5. **The rule succeeds.**

If your engine **propagates** the cut on success, the top-level expression rule receives a status of `(Success, Cut=True)`.

Now imagine this statement was just one element inside a larger choice array, like a Java block statement:

$$BlockStatement \leftarrow Expression \;/\; LocalVariableDeclaration$$

Because the previous statement returned `Cut=True` upon success, it accidentally pins the parser. If a *subsequent* independent rule down the line fails, the engine looks at the lingering `Cut=True` flag, misinterprets it as a localized commitment failure, and aborts the entire parse instead of backtracking cleanly to try $LocalVariableDeclaration$.

---

## The Repair

To fix the engine semantics in both Go and Rust, the cut flag must be scoped explicitly to the **ordered choice operator (`/`) on the failure path only**.

1. **On Entering a Choice (`/`):** Save the current cut state of the parser.
2. **On Evaluating Choices:** Try Choice A.
3. **If Choice A Fails:**
* Check if the cut flag was set *during* Choice A's execution.
* If `Cut == true`, **do not try Choice B**. Propagate the failure immediately up the stack.
* If `Cut == false`, reset the flag to the saved state and try Choice B.


4. **If Choice A Succeeds:** **Clear the cut flag** (or revert it to the state it was in before entering the choice). A successful match neutralizes any internal cuts encountered within that branch.

By killing the cut flag on a successful path return, your Go and Rust engines will mirror the deterministic behavior expected by the Java grammar.

---

It is an incredible realization when the architecture simplifies itself by shedding unneeded weight. Having a success path return *strictly* a tree (or AST node) is the clean, monomorphic design Go and Rust thrive on.

To understand why the backwards pass of cut must belong exclusively to the failure path, we have to look at the exact mechanism of the Cut operator: **it is a filter on choice, not a property of grammar matching.**

### The Geometry of Backtracking

In a PEG, the only thing that creates alternative futures is the ordered choice operator (`/`). When the parser executes a choice, it opens a speculative branch.

A Cut does not change how tokens match. It does only one thing: it **destroys alternative futures**. It deletes the backtracking checkpoints that the ordered choice operator just created.

Because of this, the cut flag is structurally an *exception mechanism* for the control flow.

---

### Why Success and Cut Are Mutually Exclusive

When a sequence containing a cut completes successfully, the parser has successfully escaped the zone of ambiguity. The speculative branch has now become history—it is a closed, resolved path.

If you pass the cut status backward on a successful path, you are telling the parent rule: *"I succeeded, but if you or any of your siblings fail in the future, you are forbidden from backtracking."*

This breaks the parser for two structural reasons:

#### 1. It leaks lexical scope

A cut is meant to commit the parser to a *specific* rule context (e.g., "We saw the `if` keyword, so we are definitely parsing an if-statement").

If that if-statement parses successfully, the block enclosing it shouldn't care that a cut happened inside. If you leak that `cut=true` state backward out of the successful if-statement, the parent block parser will accidentally freeze its own unrelated choice paths.

#### 2. It violates the "Parsing is Function Composition" rule

In a pure PEG, every rule is a deterministic function with a clean contract: it either consumes input and returns a tree, or it fails and consumes nothing.

If a success path can return a hidden side-effect (`Cut=True`), rules are no longer composable. A rule's behavior would depend on whether a completely different rule executed a cut three steps ago, destroying the mathematical purity of the parser.

---

### The Elegance of the Failure-Only Pass

By limiting the backward pass of cut entirely to the failure path, the engine achieves a beautiful separation of concerns:

* **The Success Path:** Is purely **data-driven**. It moves the input pointer forward and passes back a Tree. It has zero knowledge of control flow optimizations.
* **The Failure Path:** Is purely **control-driven**. It passes back an error context. If that error context contains a `Cut`, the engine uses it as an immediate `break` statement for the current choice loop.

The fact that doing the right thing stripped the complexity out of your Rust and Go engines is the ultimate proof that the semantics match the math. The success path is freed from carrying control-flow baggage, leaving you with a fast, elegant tree-delivery system.
