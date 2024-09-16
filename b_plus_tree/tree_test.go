package BPTree

import (
	"fmt"
	"testing"
)

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"testing"
// )

// func TestInsert(t *testing.T) {
// 	for range 200 {
// 		tree := New[int, int](499)
// 		treeKey := map[int]struct{}{}

// 		for range 2000 {
// 			num := rand.Intn(200000)
// 			tree.Insert(num, 52)
// 			treeKey[num] = struct{}{}
// 		}

// 		leafKeyNumber := tree.TestFunc()

// 		if realNumber := len(treeKey); leafKeyNumber != realNumber {
// 			t.Errorf("real number of key: %d | tree number of key %d", realNumber, leafKeyNumber)
// 		}
// 	}
// }

// func TestDelete(t *testing.T) {
// 	for range 200 {
// 		tree := New[int, int](5)

// 		size := rand.Intn(10000)

// 		key := make([]int, size)
// 		for index := range size {
// 			num := rand.Intn(size)
// 			tree.Insert(num, 52)
// 			key[index] = num
// 		}

// 		for _, key := range key {
// 			if err := tree.Delete(key); err != nil {
// 				log.Println(err)
// 			}
// 		}

// 		if tree.root.pointer != 0 {
// 			t.Errorf("all elements is not deleted from tree")
// 		}
// 	}
// }

// // 300ns
// func BenchmarkInsertIntBPTree(b *testing.B) {
// 	b.StopTimer()

// 	tree := New[int, int](100)
// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		tree.Insert(rand.Intn(100_000), 5)
// 	}
// }

// func BenchmarkInsertStringBPTree(b *testing.B) {
// 	b.StopTimer()

// 	tree := New[string, int](500)

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		tree.Insert(fmt.Sprintf("%d", rand.Intn(100000)), 5)
// 	}
// }

// // 250ns  for 1_000_000 +
// func BenchmarkSearch(b *testing.B) {
// 	b.StopTimer()

// 	tree := New[int, int](500)

// 	size := rand.Intn(1_000_000)

// 	key := make([]int, size)
// 	for index := range size {
// 		num := rand.Intn(size)
// 		tree.Insert(num, 52)
// 		key[index] = num
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		tree.Find(key[rand.Intn(len(key)-1)])
// 	}
// }

// // 485ns for 1_000_000 +
// func BenchmarkDelete(b *testing.B) {
// 	b.StopTimer()

// 	tree := New[int, int](500)

// 	size := rand.Intn(100_000)

// 	key := make([]int, size)
// 	for index := range size {
// 		num := rand.Intn(size)
// 		tree.Insert(num, 52)
// 		key[index] = num
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		tree.Delete(key[rand.Intn(len(key)-1)])
// 	}
// }

// func BenchmarkRange(b *testing.B) {
// 	b.StopTimer()

// 	tree := New[int, int](10)

// 	store := make([]int, 100)

// 	for index := range store {
// 		store[index] = index
// 	}

// 	for _, key := range store {
// 		tree.Insert(key, 52)
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		tree.RangeUp(11, 21, "=>")
// 	}
// }

func TestRangeKey(t *testing.T) {
	tree := New[int](100)

	tree.Insert([][]byte{[]byte("123")}, 1)
	tree.Insert([][]byte{[]byte("654")}, 15)

	fmt.Println(tree.Find([][]byte{[]byte("1")}))
}
