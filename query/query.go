package query

//query definition

func CreateDataBase() []byte {
	return []byte("CREATE_DATABASE Movies")
}

//KEY 		   | DATABASE | TABLE_NAME | COLUMNS AND TYPES
//CREATE_TABLE   database   TABLE_NAME  (name TEXT Age INT)

func CreateTable() []byte {
	return []byte("CREATE_TABLE Movies.Actor id UUID, name TEXT, age INT [id, name]")
}

// INSERT
// INSERT database table 5janko221
func Insert() []byte {
	return []byte("INSERT Movies.Actor 00016123456789101112100011Maze runner0000220\n")
}

// SEARCH Database Table
// name (TEXT) == "Janko" age (INT) >= 18
func Search() []byte {
	return []byte("SEARCH Movies.Actor id (UUID) == 874123 name (TEXT) == Janko")
}

func Insert2() []byte {
	return []byte("INSERT Movies.Actor 00016123456719101112100011More fnnner0000220\n")
}

func NewIndex() []byte {
	return []byte("INDEX Movies.Actor [age]")
}
