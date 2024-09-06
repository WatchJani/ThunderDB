package main

import (
	"fmt"
	"root/builder"
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

	buf := make([]byte, 4096)

	builder := builder.New(buf)

	dataInsert := []byte("Janko")
	builder.Insert(dataInsert, builder.Reservations(dataInsert))

	dataInsert = []byte("Kondic")
	builder.Insert(dataInsert, builder.Reservations(dataInsert))

	fmt.Println(buf[:15])
}
