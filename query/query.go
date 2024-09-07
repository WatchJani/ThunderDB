package query

//query definition

func CreateDataBase() []byte {
	return []byte("CREATE_DATABASE Movies")
}

//KEY 		   | DATABASE | TABLE_NAME | COLUMNS AND TYPES
//CREATE_TABLE   database   TABLE_NAME  (name TEXT Age INT)

func CreateTable() []byte {
	return []byte("CREATE_TABLE Movies.Actors id UUID, name TEXT, age INT [id, name]")
}

// INSERT
// INSERT database table 5janko221
func Insert() []byte {
	return []byte("INSERT Movies.Actor 16123456789101112111Maze runner220")
}

// SEARCH Database Table
// name (TEXT) == "Janko" age (INT) >= 18
func Search() []byte {
	return []byte("SEARCH Movies.Actor name(TEXT) == Janko age(INT) >= 18")
}
