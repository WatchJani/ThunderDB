package query

import "testing"

func TestFilter(t *testing.T) {
	// just 45 is not work
	req := []byte("45")
	operation := []byte("<")

	if ok := GenerateFilter(operation, req, "INT")([]byte("300")); !ok {
		t.Errorf("wrong operation")
	}
}

// 44ns just for creating filter function
func BenchmarkSpeedFilterGenerateFunction(b *testing.B) {
	req := []byte("45")
	operation := []byte("<")

	start := []byte("300")

	for i := 0; i < b.N; i++ {
		GenerateFilter(operation, req, "FLOAT")(start)
	}
}

func BenchmarkConverting(b *testing.B) {
	b.StopTimer()

	dr := []byte("23")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		bytesToInt(dr)
	}
}
