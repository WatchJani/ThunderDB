package database

import "root/index"

type Database struct {
	table map[string]*index.IndexBuilder //index builder is table
}

func New(name string) *Database {
	return &Database{
		table: make(map[string]*index.IndexBuilder),
	}
}

// name of table, and all information about table
func (db *Database) CreateTable() {

}
