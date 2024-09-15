package manager

import (
	"os"
)

type Manager struct {
	store    *os.File
	memTable []byte
	flush    []byte
	
}

func New(memTable []byte, file *os.File) (*Manager, error) {
	//

	return &Manager{
		store:    file,
		memTable: memTable,
		flush:    nil,
	}, nil
}
