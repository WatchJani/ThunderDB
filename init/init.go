package init

import (
	"log"
	"root/query"
	"root/thunder"
	"time"
)

func init() {
	thunder, err := thunder.New()
	if err != nil {
		log.Println(err)
		return
	}

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

	if _, err := thunder.QueryParser(query.Insert2()); err != nil {
		log.Println(err)
		return
	}
	for range 171197 {
		if _, err := thunder.QueryParser(query.Insert2()); err != nil {
			log.Println(err)
			return
		}
	}

	// if _, err := thunder.QueryParser(query.Search()); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	time.Sleep(5 * time.Second)
}
