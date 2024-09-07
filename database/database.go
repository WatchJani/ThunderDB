package database

import "root/index"

type Database struct {
	table map[string]*index.Table //index builder is table
}

func New(name string) *Database {
	return &Database{
		table: make(map[string]*index.Table),
	}
}

// name of table, and all information about table
func (db *Database) CreateTable() {

}
