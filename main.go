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

	if _, err := thunder.QueryParser(query.NewIndex()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.Insert()); err != nil {
		log.Println(err)
		return
	}

	if _, err := thunder.QueryParser(query.Search()); err != nil {
		log.Println(err)
		return
	}

	time.Sleep(5 * time.Second)
}
