package main

import (
	"fmt"
	"log"

	// _ "root/init"

	b "root/b_plus_tree"
)

func main() {
	// tree := b.New[int](100)

	// key := [][]byte{[]byte("123"), []byte("451")}
	// key2 := [][]byte{[]byte("123"), []byte("116")}

	// tree.Insert(key, 123)
	// tree.Insert(key2, 145)

	// node, index, err := tree.Find([][]byte{[]byte("123"), []byte("500")}, "<")
	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(node.GetItem(index))

	// tree := skip_list.New(32, 4000, 0.25)

	// key := [][]byte{[]byte("1"), []byte("34")}
	// key2 := [][]byte{[]byte("1"), []byte("44")}
	// key3 := [][]byte{[]byte("1"), []byte("74")}

	// key4 := [][]byte{[]byte("2"), []byte("45")}
	// key5 := [][]byte{[]byte("2"), []byte("55")}

	// tree.Insert(key3, 74)
	// tree.Insert(key, 34)
	// tree.Insert(key5, 55)
	// tree.Insert(key2, 44)
	// tree.Insert(key4, 45)

	// tree.ReadAllFromLeftToRight()

	// ok, offset := tree.Search([][]byte{[]byte("1"), []byte("35")}, "<")
	// if !ok {
	// 	log.Println("not found")
	// }

	// fmt.Println(offset)

	tree := b.New[int](11)

	key := [][]byte{[]byte("1111")}
	tree.Insert(key, 1)
	tree.Insert(key, 2)
	tree.Insert(key, 3)
	tree.Insert(key, 4)
	tree.Insert(key, 5)
	tree.Insert(key, 6)
	tree.Insert(key, 7)
	tree.Insert(key, 8)
	tree.Insert(key, 9)
	tree.Insert(key, 10)

	node, index, err := tree.Find(key, "==")
	fmt.Println("index", index)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(node.GetValue(index))

	fmt.Println("test")
	tree.TestRoot()
}
