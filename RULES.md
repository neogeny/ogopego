---
apply: always
---

# Rules

## Interaction

- Never change files without a plan and user authorization
- Never use a Git command that alters files or their version control status.
- Never change more than the files explicitly named in the authorization
- Always consult with the User before making changes that impact multiple files
- Evaluate changes before applying them (no "apply and see" approach)
- Do not act on assumptions. Always verify assumptions with the User.
- Always read a file again to verify its current state before making changes. Never assume state from memory.


## Justfile

- There is a `Justfile` defined with targets with the most common tasks. For
  consistency it is preferred to invoke the `just` target instead of the
  direct command on the command line.

## Testing

- Tests can be marked to be skipped/ignored, but must compile
    ```

## Code Tools

- Do not use `sed`, `awk`, or similar text tools to modify code
- When using `ripgrep` `rg`, explicitly name the directories to search (e.g.,
  `rg -l 'pattern' src tests`) instead of excluding directories
- Create a `./tmp/` directory for temporary files instead of using `/tmp`
- Never try to access or modify a file or directory outside the current
  project's directory. If there's additional files relevant to the project's
  context they will be provided through a symlink.
