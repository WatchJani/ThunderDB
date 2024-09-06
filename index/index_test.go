package index

import (
	"log"
	"root/column"
	"testing"
)

// 6ns
func BenchmarkChoice(b *testing.B) {
	b.StopTimer()

	columns := []column.Column{
		column.New("id", "UUID"),
		column.New("name", "TEXT"),
		column.New("phone", "TEXT"),
		column.New("sex", "TEXT"),
		column.New("age", "int"),
	}

	clusterIndex := New("cluster", "id", "age")

	IndexBuilder, err := NewIndexBuilder(columns, clusterIndex)
	if err != nil {
		log.Println(err)
		return
	}
	IndexBuilder.AddIndex(New("age", "age"))

	userField := []string{"id"}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		IndexBuilder.Choice(userField)
	}
}
