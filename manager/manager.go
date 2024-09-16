package manager

import (
	"os"
)

type Manager struct {
	store    *os.File
	memTable []byte
	old      []byte
}

func (m *Manager) SetOld(data []byte) {
	m.old = data
}

func (m *Manager) GetOld() *[]byte {
	return &m.old
}

func (m *Manager) SetGetOld(data []byte) *[]byte {
	m.old = data
	return &m.old
}

func New(memTable []byte, file *os.File) (*Manager, error) {
	//

	return &Manager{
		store:    file,
		memTable: memTable,
		old:      nil,
	}, nil
}
