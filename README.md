[![Go Reference](https://pkg.go.dev/badge/fortio.org/sets.svg)](https://pkg.go.dev/fortio.org/sets)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/sets)](https://goreportcard.com/report/fortio.org/sets)
[![GitHub Release](https://img.shields.io/github/release/fortio/sets.svg?style=flat)](https://github.com/fortio/sets/releases/)
[![Coverage](https://codecov.io/github/fortio/sets/branch/main/graph/badge.svg?token=LONYZDFQ7C)](https://codecov.io/github/fortio/sets)

# sets
Sets and Set operations in golang 

Internally as simple as a `map[T]struct{}` but with additional helper functions, operations like Union, Intersect, XOR, to/from Slices conversion, etc...
and hiding the `struct{}{}` etc...
