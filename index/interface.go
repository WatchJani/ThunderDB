package index

type Index interface {
	Search()
	Insert([][]byte, int)
	GetByColumn(int) string
	GetColumnNumber() int
}
