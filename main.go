package main

import (
	"fmt"
	"log"
	"root/index"
)

func main() {
	columns := []string{"id", "name", "phone", "sex", "age"}

	clusterIndex := index.New("cluster", []string{"id", "age"})

	IndexBuilder, err := index.NewIndexBuilder(columns, []index.Index{clusterIndex})
	if err != nil {
		log.Println(err)
	}

	userField := []string{"id"}

	fmt.Println(IndexBuilder.Choice(userField))
}
