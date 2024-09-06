package query

import (
	"bytes"
)

type Filter func([]byte) bool

func GenerateFilter(operation, input []byte) func([]byte) bool {
	switch {
	case bytes.Equal(operation, []byte("<=")):
		return LessThanOrEqualTo(input)
	case bytes.Equal(operation, []byte("<")):
		return LessThan(input)
	case bytes.Equal(operation, []byte(">")):
		return GreaterThan(input)
	case bytes.Equal(operation, []byte(">=")):
		return GreaterOrEqualTo(input)
	default:
		return EqualTo(input)
	}
}

// [<=]
func LessThanOrEqualTo(input []byte) func([]byte) bool {
	return func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 0 || result == -1
	}
}

// [<]
func LessThan(input []byte) func([]byte) bool {
	return func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == -1
	}
}

// [>]
func GreaterThan(input []byte) func([]byte) bool {
	return func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 1
	}
}

// [>=]
func GreaterOrEqualTo(input []byte) func([]byte) bool {
	return func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 0 || result == 1
	}
}

// [==]
func EqualTo(input []byte) func([]byte) bool {
	return func(b []byte) bool {
		return bytes.Equal(input, b)
	}
}
