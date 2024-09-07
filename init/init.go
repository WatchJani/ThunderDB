package init

import (
	"root/query"
	"root/thunder"
)

func init() {
	thunder := thunder.New()

	createDatabase := query.CreateDataBase()
	thunder.QueryParser(createDatabase)

	createTable := query.CreateTable()
	thunder.QueryParser(createTable)
	
}
