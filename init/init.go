package init

import (
	"log"
	"root/query"
	"root/thunder"
)

func init() {
	thunder := thunder.New()

	createDatabase := query.CreateDataBase()
	if err := thunder.QueryParser(createDatabase); err != nil {
		log.Println(err)
	}

	createTable := query.CreateTable()
	if err := thunder.QueryParser(createTable); err != nil {
		log.Println(err)
	}

	insertQuery := query.Insert()
	if err := thunder.QueryParser(insertQuery); err != nil {
		log.Println(err)
	}

	insertQuery = query.Insert2()
	if err := thunder.QueryParser(insertQuery); err != nil {
		log.Println(err)
	}

	// searchQuery := query.Search()
	// if err := thunder.QueryParser(searchQuery); err != nil {
	// 	log.Println(err)
	// }
}
