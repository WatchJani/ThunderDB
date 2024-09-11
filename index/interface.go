package index

type Index interface {
	Search()
	Insert([][]byte)
	GetByColumn(int) string
	GetColumnNumber() int
}
