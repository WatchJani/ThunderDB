package manager

import (
	"os"
	"root/skip_list"

	"github.com/WatchJani/stack"
)

type Manager struct {
	store    *os.File
	memTable []byte

	old  []byte
	tree *skip_list.SkipList

	stack stack.Stack[[]byte]
}

func (m *Manager) GetAllData() (*os.File, []byte, []byte, *skip_list.SkipList) {
	return m.store, m.memTable, m.old, m.tree
}

func (m *Manager) SetOld(data []byte, tree *skip_list.SkipList) {
	m.old = data
	m.tree = tree
}

func (m *Manager) GetOld() *[]byte {
	return &m.old
}

func New(memTable []byte, file *os.File) (*Manager, error) {

	size := 10
	stack := stack.New[[]byte](size)
	for range size {
		stack.Push(make([]byte, 4096))
	}

	//
	return &Manager{
		store:    file,
		memTable: memTable,
		old:      nil,
		tree:     nil,
		stack:    stack,
	}, nil
}

func (m *Manager) GetFreeByte() []byte {
	freeByte, err := m.stack.Pop()
	if err != nil {
		freeByte = make([]byte, 4096)
	}

	return freeByte
}

func (m *Manager) FlushFreeByte(freeByte []byte) {
	m.stack.Push(freeByte)
}
