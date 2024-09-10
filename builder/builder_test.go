package builder

// func TestInsert(t *testing.T) {
// 	buf := make([]byte, 4096)

// 	builder := New(buf)

// 	// dataInsert := []byte("Janko")
// 	// builder.ParallelWrite(dataInsert)

// 	// dataInsert = []byte("Kondic")
// 	// builder.ParallelWrite(dataInsert)

// 	result := []byte("JankoKondic")

// 	if res := bytes.Compare(buf[:builder.counter], result); res != 0 {
// 		t.Errorf("%v != %v ", buf[:builder.counter], result)
// 	}
// }

// // 9ns
// func BenchmarkInsert(b *testing.B) {
// 	b.StopTimer()
// 	buf := make([]byte, 4096)

// 	builder := New(buf)

// 	// dataInsert := []byte("Janko")

// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		// builder.ParallelWrite(dataInsert)
// 		builder.Reset()
// 	}
// }
