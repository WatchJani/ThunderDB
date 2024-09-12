package main

import (
	"fmt"
	_ "root/init"

	t "root/b_plus_tree"
)

func main() {
	// columns := []column.Column{
	// 	column.New("id", "UUID"),
	// 	column.New("name", "TEXT"),
	// 	column.New("phone", "TEXT"),
	// 	column.New("sex", "TEXT"),
	// 	column.New("age", "int"),
	// }

	// clusterIndex := index.New("cluster", "id", "age")

	// IndexBuilder, err := index.NewIndexBuilder(columns, clusterIndex)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// IndexBuilder.AddIndex(index.New("age", "age"))

	// userField := []string{"age"}

	// fmt.Println(IndexBuilder.Choice(userField))

	//)====================================================0
	// buf := make([]byte, 4096)

	// builder := builder.New(buf)

	// dataInsert := []byte("Janko")
	// builder.ParallelWrite(dataInsert)

	// dataInsert = []byte("Kondic")
	// builder.ParallelWrite(dataInsert)

	// fmt.Println(buf[:15])

	//====================================================

	tree := t.New[int](5)

	tree.Insert([][]byte{[]byte("005")}, 1)
	tree.Insert([][]byte{[]byte("100")}, 2)
	tree.Insert([][]byte{[]byte("153")}, 3)
	tree.Insert([][]byte{[]byte("251")}, 4)
	tree.Insert([][]byte{[]byte("357")}, 5)
	tree.Insert([][]byte{[]byte("654")}, 6)

	fmt.Println(tree.BetweenKey([][]byte{[]byte("099")}))
}
