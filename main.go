package main

import (

	// _ "root/init"

	"fmt"
	"log"
	"root/skip_list"
)

func main() {
	// tree := b.New[int](100)

	// key1 := [][]byte{[]byte("123"), []byte("451")}
	// key2 := [][]byte{[]byte("123"), []byte("116")}

	// tree.Insert(key1, 123)
	// tree.Insert(key2, 145)

	// node, index, err := tree.Find([][]byte{[]byte("123"), []byte("500")}, "<")
	// if err != nil {
	// 	log.Println(err)
	// }

	// fmt.Println(node.GetItem(index))

	tree := skip_list.New(32, 4000, 0.25)

	key1 := [][]byte{[]byte("1"), []byte("34")}
	key2 := [][]byte{[]byte("1"), []byte("44")}
	key3 := [][]byte{[]byte("1"), []byte("74")}

	key4 := [][]byte{[]byte("2"), []byte("45")}
	key5 := [][]byte{[]byte("2"), []byte("55")}

	tree.Insert(key3, 74)
	tree.Insert(key1, 34)
	tree.Insert(key5, 55)
	tree.Insert(key2, 44)
	tree.Insert(key4, 45)


	tree.ReadAllFromLeftToRight()

	ok, offset := tree.Search([][]byte{[]byte("2"), []byte("56")})
	if !ok {
		log.Println("not found")
	}

	fmt.Println(offset)
}
