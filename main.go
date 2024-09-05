package main

import (
	"fmt"
	"log"
	"root/column"
	"root/index"
)

func main() {
	columns := []column.Column{
		column.New("id", "UUID"),
		column.New("name", "TEXT"),
		column.New("phone", "TEXT"),
		column.New("sex", "TEXT"),
		column.New("age", "int"),
	}

	clusterIndex := index.New("cluster", "id", "age")

	IndexBuilder, err := index.NewIndexBuilder(columns, clusterIndex)
	if err != nil {
		log.Println(err)
		return
	}
	IndexBuilder.AddIndex(index.New("age", "age"))

	userField := []string{"age"}

	fmt.Println(IndexBuilder.Choice(userField))
}
