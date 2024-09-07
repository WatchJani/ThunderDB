package query

//query definition

//INSERT Database Table
//name (TEXT) == "Janko" age (INT) >= 18
//skip

//KEY 		   | DATABASE | TABLE_NAME | COLUMNS AND TYPES
//CREATE_TABLE   database   TABLE_NAME  (name TEXT Age INT)

func CreateDataBase() []byte {
	return []byte("CREATE_DATABASE Movies")
}

func CreateTable() []byte {
	return []byte("CREATE_TABLE Movies.Actors id UUID, name TEXT, age INT [id, name]")
}
