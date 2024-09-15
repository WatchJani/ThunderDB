package thunder

import (
	"log"
	"root/cutter"
	"root/database"
	"root/linker"
)

type Thunder struct {
	linker.Linker
	database map[string]*database.Database
}

func New() Thunder {
	linker := linker.New()

	cutter, err := cutter.New(linker, "/home/janko/Desktop/chanel23l/store.bin", 10)
	if err != nil {
		log.Println(err)
	}

	go cutter.Cut()

	return Thunder{
		Linker:   linker,
		database: make(map[string]*database.Database),
	}
}

func (t *Thunder) NewDatabase(name string) error {
	t.database[name] = database.New(t.Linker)
	return nil
}

func (t *Thunder) NewTable(databaseName, tableName string) error {
	database := t.database[databaseName]
	database.CreateTable(tableName)

	return nil
}

func (t *Thunder) Inset(databaseName, tableName string, data []byte) {
	database := t.database[databaseName]
	table := database.GetTable(tableName)

	table.Insert(data)
}
