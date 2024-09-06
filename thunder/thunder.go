package thunder

import (
	"fmt"
	"root/database"
)

type Thunder struct {
	thunder map[string]*database.Database
}

func New() *Thunder {
	return &Thunder{
		thunder: make(map[string]*database.Database),
	}
}

func (t *Thunder) CreateDatabase(name string) (*database.Database, error) {
	if database, ok := t.thunder[name]; ok {
		return database, fmt.Errorf("database [%s] is exist", name)
	}

	return database.New(name), nil
}
