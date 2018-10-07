# JEX

JEX is a proof-of-concept of adding exceptions to Go programs.

This is an AST-based code generator with a tiny runtime dependency. 
Calls to special macros `THROW` and `TRY` are replaced with appropriate combinations of panics, defers and recovers.
Also it provides convenient way to check `error` values by introduction of special macro-variable `ERR`.

## How to enable

Add to the beginning of your file:

```go
//+build jex
//go:generate jex
```

Add dot-import

```go
import . "github.com/anjensan/jex"
```

Now you can run `go generate -tags jex` to produce exceptionally new code.

## How to use

Throw an exception

```go
THROW(errors.New("error name"))
```

Attach debug info (variable, stacktrace etc)

```go
THROW(errors.New("error name"), "foo", "bar", debug.Stack())
```

Catch and handle

```go
if TRY() {
    // some code

} else {
    // EX() expanded into local variable name 
    // EX().Err() - untouched original `error`
    fmt.Println(EX().Err())

    // Print extra info
    for _, l := range EX().Logs() { fmt.Println(l) }

    // Suppressed errors
    for _, ex := range EX().Suppressed() { fmt.Println(ex.Err()) }

    // Print *correct* stacktrace, showing where original panic was called!
    debug.PrintStack()

    // Print all info at the same time
    fmt.Println(EX())
}
```

Rethrow an exception

```go
if TRY() {
    // some code

} else {
    // Throw new exception keeping the current one as suppressed
    THROW(errors.New("new error"))
}
```

Integration with old-style functions

```go
// jex adds `if ERR != nil { THROW(err) }` after such assignments
file, ERR := os.Open(filename)
```

All exceptional functions must be suffixed with `_`

```go
func printSum_(a, b string) {
    x, ERR := strconv.Atoi(a)
    y, ERR := strconv.Atoi(b)
    fmt.Println("result:", x + y)
}
```

It is allowed to use `return` & `defer` inside `if-TRY` blocks

```go
func readFile(fn string) int {
    if TRY() {
        r, ERR := os.Open(src)
        defer r.Close()
        return parseReader_(r)
    } else {
        log.Errorf("Exception %v", EX())
        return 0
    }
}
```

See file `demo/demo.go` for more examples.
