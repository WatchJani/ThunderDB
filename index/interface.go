package index

type Index interface {
	Insert([][]byte, int)
	GetByColumn() []string
	GetIndexType() string
}
