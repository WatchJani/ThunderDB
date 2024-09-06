package query

import (
	"bytes"
	"strconv"
)

type Filter func([]byte) bool

func GenerateFilter(operation, input []byte, dataType string) func([]byte) bool {
	switch {
	case bytes.Equal(operation, []byte("<=")):
		return LessThanOrEqualTo(input, dataType)
	case bytes.Equal(operation, []byte("<")):
		return LessThan(input, dataType)
	case bytes.Equal(operation, []byte(">")):
		return GreaterThan(input, dataType)
	case bytes.Equal(operation, []byte(">=")):
		return GreaterOrEqualTo(input, dataType)
	default:
		return EqualTo(input, dataType)
	}
}

func bytesToInt(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}

func bytesToFloat(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}

func Response(input []byte, dataType string, Lexicographic func([]byte) bool) func([]byte) bool {
	switch dataType {
	case "INT":
		ClientQuery, _ := bytesToInt(input)
		return func(b []byte) bool {
			numberFromStore, _ := bytesToInt(b)
			return ClientQuery <= numberFromStore
		}
	case "FLOAT":
		ClientQuery, _ := bytesToFloat(input)
		return func(b []byte) bool {
			numberFromStore, _ := bytesToFloat(b)
			return ClientQuery <= numberFromStore
		}
	default:
		return Lexicographic
	}
}

// [<=]
func LessThanOrEqualTo(input []byte, dataType string) func([]byte) bool {
	return Response(input, dataType, func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 0 || result == -1
	})
}

// [<]
func LessThan(input []byte, dataType string) func([]byte) bool {
	return Response(input, dataType, func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == -1
	})
}

// [>]
func GreaterThan(input []byte, dataType string) func([]byte) bool {
	return Response(input, dataType, func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 1
	})
}

// [>=]
func GreaterOrEqualTo(input []byte, dataType string) func([]byte) bool {
	return Response(input, dataType, func(b []byte) bool {
		result := bytes.Compare(input, b)
		return result == 0 || result == 1
	})
}

// [==]
func EqualTo(input []byte, dataType string) func([]byte) bool {
	return Response(input, dataType, func(b []byte) bool {
		return bytes.Equal(input, b)
	})
}
