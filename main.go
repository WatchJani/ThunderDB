package main

import (
	"log"
	"root/query"
	"root/thunder"
	"time"
)

func main() {
	thunder := thunder.New()

	if _, err := thunder.QueryParser(query.CreateDataBase()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.CreateTable()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.Insert()); err != nil {
		log.Println(err)
		return
	}

	// if err := thunder.NewTable("netflix", "user", []column.Column{}); err != nil {
	// 	log.Println(err)
	// }

	// for range 171197 * 2 {
	// 	thunder.Inset("netflix", "user", query.Insert()[20:])
	// }

	time.Sleep(5 * time.Second)
}
