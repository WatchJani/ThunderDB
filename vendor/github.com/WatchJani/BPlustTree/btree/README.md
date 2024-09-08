# Go B+ Tree Library

[![GoDoc](https://godoc.org/github.com/WatchJani/BPlustTree?status.svg)](https://pkg.go.dev/github.com/yourusername/bplustree)
[![Go Report Card](https://goreportcard.com/badge/github.com/WatchJani/BPlustTree)](https://goreportcard.com/report/github.com/yourusername/bplustree)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A high-performance, thread-safe implementation of the B+ Tree data structure in Go. Ideal for database indexing, sorted data management, and efficient range queries.

## Features

- **Full B+ Tree Implementation**: Supports insertion, deletion, and search operations.
- **Range Queries**: Efficiently handles range queries for sorted data.
- **Customizable Node Size**: Users can define the maximum number of keys per node.
- **Iterator Support**: Provides a built-in iterator for in-order traversal.
- **Thread-Safe**: Designed for concurrent read and write operations.
- **Optimized for Performance**: Tuned for high performance in various use cases.

## Installation

To install the package, use `go get`:

```sh
go get github.com/WatchJani/BPlustTree
```
### Explanation:

1. **Go Code Block**: The Go code is enclosed in triple backticks with `go` specified for syntax highlighting:

```go
package main

import (
    "fmt"
    "log"

    t "github.com/WatchJani/BPlustTree"
)

func main() {
	BPTree := t.New[int, int](50)

	BPTree.Insert(123, 123)

	value, err := BPTree.Find(123)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(value)

	if err := BPTree.Delete(123); err != nil {
		log.Println(err)
	}
}

```

2. **Markdown Sections**: Standard Markdown syntax is used for headings (`#`, `##`), lists, and other text formatting.

3. **Shell Commands**: For shell commands like `go get`, use triple backticks with `sh` for shell syntax highlighting:
    ```sh
    go get github.com/WatchJani/BPlustTree
    ```

### Output:
When rendered on GitHub, the example Go code will be syntax-highlighted, making it easier to read and understand. This approach is useful for providing clear and well-documented examples in your project's README file.
