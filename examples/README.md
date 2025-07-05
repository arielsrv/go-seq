# go-seq

[![Go Reference](https://pkg.go.dev/badge/github.com/arielsrv/go-seq.svg)](https://pkg.go.dev/github.com/arielsrv/go-seq)
[![codecov](https://codecov.io/gh/arielsrv/go-seq/branch/main/graph/badge.svg)](https://codecov.io/gh/arielsrv/go-seq)

A functional, LINQ-inspired, type-safe sequence library for Go.

## Features

- Functional-style sequence operations (Where, Select, SelectMany, etc.)
- Type-safe, generic API
- Works with any type, not just slices
- Lazy evaluation
- Easily convert sequences to slices or maps

## Installation

```sh
go get github.com/arielsrv/go-seq
```

## Usage

```go
import (
    "fmt"
    "strings"
    "github.com/arielsrv/go-seq"
    "github.com/arielsrv/go-seq/iter"
)

func main() {
    people := []string{"John,Engineer,25", "Mary,Doctor,30"}
    result := seq.Select(
        seq.SelectMany(
            seq.Where(
                seq.Yield(people...),
                func(p string) bool { return strings.Contains(strings.ToLower(p), "engineer") },
            ),
            func(p string) iter.Seq[string] {
                parts := strings.Split(p, ",")
                return seq.Yield(parts[0] + "-" + parts[1])
            },
        ),
        strings.ToUpper,
    )
    fmt.Println(seq.ToSlice(result)) // Output: [JOHN-ENGINEER]
}
```

See more examples in [`examples/full/main.go`](examples/full/main.go).

## API Highlights

- `Where(seq, predicate)` – filter values
- `Select(seq, mapper)` – map values
- `SelectMany(seq, mapper)` – flatMap
- `ToSlice(seq)` – convert to slice
- `ToMap(seq2)` – convert to map

## Running Tests

```sh
go test -v -cover ./...
```

## License

MIT 