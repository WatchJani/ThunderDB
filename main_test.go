package main

import (
	"log"
	"root/column"
	"root/query"
	"root/thunder"
	"testing"
)

func BenchmarkInsertSpeed(b *testing.B) {
	b.StopTimer()

	thunder := thunder.New()
	if err := thunder.NewDatabase("netflix"); err != nil {
		log.Println(err)
	}

	if err := thunder.NewTable("netflix", "user", []column.Column{}); err != nil {
		log.Println(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		thunder.Inset("netflix", "user", query.Insert()[20:])
	}
}
