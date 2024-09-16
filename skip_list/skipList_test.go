package skip_list

import (
	"testing"
)

func TestInsert(t *testing.T) {
	skipList := New(32, 2000, 0.25)

	// user_id	timestamp	event_type
	key1 := [][]byte{[]byte("1"), []byte("2024-09-01 12:34:56"), []byte("login")}
	skipList.Insert(key1, 5)

	key2 := [][]byte{[]byte("2"), []byte("2024-09-01 12:35:00"), []byte("purchase")}
	skipList.Insert(key2, 2)

	key3 := [][]byte{[]byte("3"), []byte("2024-09-01 12:36:22"), []byte("logout")}
	skipList.Insert(key3, 3)

	key := [][]byte{[]byte("3")}
	if ok, value := skipList.Search(key); !ok || value != 3 {
		t.Errorf("key is not founded %v", key)
	}
}

func Benchmark(b *testing.B) {
	b.StopTimer()
	skipList := New(32, 2000, 0.25)
	key := [][]byte{[]byte("3")}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		skipList.Insert(key, 2)
	}
}
