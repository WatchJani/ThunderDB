package query

func Insert() []byte {
	return []byte("INSERT Movies Actor 0004400016123456789101112100011Maze runner0000220")
}

func CreateDataBase() []byte {
	return []byte("CREATE_DATABASE Movies")
}

// by default have id UUID
func CreateTable() []byte {
	return []byte("CREATE_TABLE Movies.Actor name TEXT age INT")
}

func NewIndex() []byte {
	return []byte("INDEX Movies.Actor name age")
}

func Search() []byte {
	return []byte("SEARCH Movies Actor name < 4511 age < 34")
}
