package main

import (
	"log"
	"root/query"
	"root/thunder"
	"testing"
)

func BenchmarkInsertSpeed(b *testing.B) {
	b.StopTimer()

	thunder := thunder.New()

	if _, err := thunder.QueryParser(query.CreateDataBase()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.CreateTable()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.NewIndex()); err != nil {
		log.Println(err)
		return
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if _, err := thunder.QueryParser(query.Insert()); err != nil {
			log.Println(err)
			return
		}
	}
}

func BenchmarkSearchSpeed(b *testing.B) {
	b.StopTimer()

	thunder := thunder.New()

	if _, err := thunder.QueryParser(query.CreateDataBase()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.CreateTable()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.NewIndex()); err != nil {
		log.Println(err)
		return
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if _, err := thunder.QueryParser(query.Search()); err != nil {
			log.Println(err)
			return
		}
	}
}
