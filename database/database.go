package database

import (
	"log"
	"os"
	"root/linker"
	"root/table"
)

type Database struct {
	linker.Linker
	table map[string]*table.Table
}

func New(linker linker.Linker) *Database {
	return &Database{
		table:  make(map[string]*table.Table),
		Linker: linker,
	}
}

func (db *Database) CreateTable(tableName string) error {
	file, err := os.OpenFile("/home/janko/Desktop/chanel23l/store.bin", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}

	table, err := table.New(db.Linker, file)
	if err != nil {
		return err
	}

	db.table[tableName] = table

	return nil
}

func (db *Database) GetTable(tableName string) *table.Table {
	return db.table[tableName]
}
