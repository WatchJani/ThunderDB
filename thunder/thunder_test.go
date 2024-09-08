package thunder

import (
	"log"
	"root/query"
	"testing"
)

func BenchmarkSpeedReqHandle(b *testing.B) {
	b.StopTimer()

	thunder := New()

	createDatabase := query.CreateDataBase()
	if err := thunder.QueryParser(createDatabase); err != nil {
		log.Println(err)
	}

	createTable := query.CreateTable()
	if err := thunder.QueryParser(createTable); err != nil {
		log.Println(err)
	}

	searchQuery := query.Search()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if err := thunder.QueryParser(searchQuery); err != nil {
			log.Println(err)
		}
	}
}
