package query

import "testing"

func TestFilter(t *testing.T) {
	// just 45 is not work
	req := []byte("045")
	operation := []byte("<")

	if ok := GenerateFilter(operation, req)([]byte("300")); !ok {
		t.Errorf("wrong operation")
	}
}

// 44ns just for creating filter function
func BenchmarkSpeedFilterGenerateFunction(b *testing.B) {
	req := []byte("045")
	operation := []byte("<")

	start := []byte("300")

	for i := 0; i < b.N; i++ {
		GenerateFilter(operation, req)(start)
	}
}
