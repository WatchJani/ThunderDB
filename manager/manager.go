package manager

import (
	"os"
	"root/skip_list"
)

type Manager struct {
	store    *os.File
	memTable []byte

	old  []byte
	tree *skip_list.SkipList
}

func (m *Manager) SetOld(data []byte, tree *skip_list.SkipList) {
	m.old = data
	m.tree = tree
}

func (m *Manager) GetOld() *[]byte {
	return &m.old
}

func New(memTable []byte, file *os.File) (*Manager, error) {
	//

	return &Manager{
		store:    file,
		memTable: memTable,
		old:      nil,
		tree:     nil,
	}, nil
}
