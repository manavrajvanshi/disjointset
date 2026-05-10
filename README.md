# disjointset

A generic Disjoint Set Union (DSU) data structure for Go, implementing **union by rank** and **path compression** for near-constant time operations.

[![Test Coverage](https://codecov.io/gh/manavrajvanshi/disjointset/branch/main/graph/badge.svg)](https://codecov.io/gh/manavrajvanshi/disjointset)
![Go Version](https://img.shields.io/badge/go-1.26.3%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## What is a Disjoint Set?

A Disjoint Set Union (also known as Union-Find) is a data structure that tracks a collection of elements partitioned into non-overlapping (disjoint) subsets. It efficiently supports two operations:

- **Union** — merge two subsets into one
- **Find** — determine whether two elements belong to the same subset

This library uses two optimizations to achieve nearly O(1) amortized time per operation:

- **Union by Rank** — the root with the higher rank becomes the parent, keeping trees shallow
- **Path Compression** — flattens the tree on every `Root` call so future lookups are faster

## Installation

```bash
go get github.com/manavrajvanshi/disjointset
```

Requires Go 1.18+ (uses generics).

## Usage

```go
package main

import (
    "fmt"
    "github.com/manavrajvanshi/disjointset"
)

func main() {
    d := disjointset.NewDSU[int]()

    // Register nodes
    d.Add(1)
    d.Add(2)
    d.Add(3)
    d.Add(4)

    fmt.Println(d.Components()) // 4 — each node is its own component

    // Merge components
    d.Union(1, 2)
    d.Union(3, 4)

    fmt.Println(d.Components()) // 2 — two groups: {1,2} and {3,4}

    // Check connectivity
    connected, _ := d.Find(1, 2)
    fmt.Println(connected) // true

    connected, _ = d.Find(1, 3)
    fmt.Println(connected) // false

    // Merge the two groups
    d.Union(2, 3)
    fmt.Println(d.Components()) // 1

    // Get the representative root of a component
    root, _ := d.Root(4)
    fmt.Println(root) // 1 (the root of the merged component)
}
```

### Using with other types

Because the DSU is generic, it works with any `comparable` type:

```go
// Strings
d := disjointset.NewDSU[string]()
d.Add("alice")
d.Add("bob")
d.Union("alice", "bob")

// Custom comparable types
type NodeID int32
d2 := disjointset.NewDSU[NodeID]()
```

## API Reference

### `NewDSU[T comparable]() DSU[T]`
Creates and returns an empty DSU.

### `Add(p T) bool`
Registers a new node `p`. Returns `true` if added, `false` if already present.

### `Union(p, q T) (bool, error)`
Merges the components containing `p` and `q`. Returns `true` if they were in different components (i.e., a merge actually occurred), `false` if they were already connected. Returns an error if either node is not registered.

### `Find(p, q T) (bool, error)`
Reports whether `p` and `q` belong to the same component. Returns an error if either node is not registered.

### `Root(p T) (T, error)`
Returns the representative (root) node of the component containing `p`. Also applies path compression as a side effect. Returns an error if `p` is not registered.

### `Components() int`
Returns the current number of disjoint sets.

## Running Tests

```bash
go test ./... -v
```

To also check for race conditions:

```bash
go test ./... -race
```

## Project Structure

```
disjointset/
├── dsu.go          # DSU interface and implementation
├── dsu_test.go     # Comprehensive table-driven tests
├── go.mod
└── .github/
    └── workflows/  # GitHub Actions CI
```

## License

[MIT](LICENSE)