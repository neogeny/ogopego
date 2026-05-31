# regexp2: `MatchString` allocates and decodes `[]rune` on every call

## Summary

`Regexp.MatchString(s)` always allocates a new `[]rune` from the input string,
even when the same `Regexp` instance is called repeatedly. For a hot loop in
a PEG parser this produces ~400 GB of garbage collection pressure per session.

The library already has a `runeCache` field on `Regexp` and a
`getRunesAndStart` helper that reuses the buffer, but `MatchString` bypasses
both — it calls the plain `getRunes(s)` which is just `[]rune(s)`.

## Impact

In the ogopego PEG parser, `MatchString` is called on every pattern match
(which may happen tens of millions of times per parse). The resulting
allocation dominated the heap profile at **95.27% of allocated space** and
**48% of CPU time** (the UTF-8 decode iteration).

After adding a buffer-reuse scheme (described below), allocation dropped
21× (398 GB → 18.9 GB), but the remaining CPU is still spent iterating the
input string rune-by-rune to fill the buffer on every match call.

## Root Cause

`regexp.go` has two code paths:

```go
// Path A — used by MatchString (ALWAYS allocates)
func getRunes(s string) []rune {
    return []rune(s)  // fresh allocation every call
}

// Path B — used by FindStringMatchStartingAt (sometimes reuses buffer)
func (re *Regexp) getRunesAndStart(s string, startAt int) ([]rune, int) {
    if startAt < 0 {
        return getRunes(s), 0  // still allocates!
    }
    if cap(re.runeCache) < len(s) {
        re.runeCache = make([]rune, len(s))  // grow if needed
    }
    i := 0
    runeIdx := -1
    for strIdx, r := range s {
        if strIdx == startAt {
            runeIdx = i
        }
        re.runeCache[i] = r
        i++
    }
    if startAt == len(s) {
        runeIdx = i
    }
    return re.runeCache[:i], runeIdx
}
```

`MatchString` at line 236 calls `run(true, -1, getRunes(s))` — path A, always
allocates. `FindStringMatchStartingAt` at line 192 calls `getRunesAndStart` —
path B, reuses buffer but still iterates the string.

## Simple Workaround (what we did in our fork)

We added a `runeCache []rune` field to `Regexp` and rewrote `getRunes` to
reuse it:

```go
func (re *Regexp) getRunesAndStart(s string, startAt int) ([]rune, int) {
    // ... existing implementation with re.runeCache reuse ...
}

// New helper used by MatchString, FindStringMatch, etc.
func (re *Regexp) cachedRunes(s string) []rune {
    if cap(re.runeCache) < len(s) {
        re.runeCache = make([]rune, len(s))
    }
    i := 0
    for _, r := range s {
        re.runeCache[i] = r
        i++
    }
    return re.runeCache[:i]
}
```

Then changed `MatchString` and `FindStringMatch` to call `cachedRunes(s)`
instead of `getRunes(s)`.

## Remaining Issue

Even with buffer reuse, the `for _, r := range s` loop must decode UTF-8
from the string on every match call. For a parser that matches against
different substrings of the same input (the cursor advances through the
text), this means re-decoding a large contiguous prefix of the input on
every pattern match.

## Ideal Solution

Add public methods that accept `[]rune` directly so callers who keep a
rune-slice can bypass string→rune conversion entirely:

```go
// Already exists:
func (re *Regexp) MatchRunes(r []rune) (bool, error)
func (re *Regexp) FindRunesMatch(r []rune) (*Match, error)
func (re *Regexp) FindRunesMatchStartingAt(r []rune, startAt int) (*Match, error)
```

These methods are already implemented and pass the rune slice directly to
`run()` with zero conversion. The gap is ergonomic: callers must manage the
`[]rune` lifetime themselves.

## Recommendations

1. **Make `MatchString`/`FindStringMatch` use the `runeCache`** to avoid
   per-call allocation (the 21× improvement we saw).
2. **Add a `ReuseRunes(s string) []rune` public method** so callers can
   decode once and reuse across multiple patterns:

   ```go
   // Caller side:
   runes := re.ReuseRunes(inputText)  // decode once
   for each pattern match {
       m, _ := re.FindRunesMatchStartingAt(runes, offset)
   }
   ```

3. **Document the `FindRunesMatch`/`MatchRunes` family** prominently so
   performance-sensitive callers know to use them.

## Environment

- Library: [regexp2](https://github.com/dlclark/regexp2) v1.12.0
- Go version: 1.24
- Use case: Repeated pattern matching against advancing substrings of a
  single input (PEG parser cursor)
