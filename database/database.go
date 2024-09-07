package database

import (
	"fmt"
	"root/column"
	"root/index"
)

type Database struct {
	table map[string]*index.Table //index builder is table
}

func New(name string) *Database {
	return &Database{
		table: make(map[string]*index.Table),
	}
}

// name of table, and all information about table
func (db *Database) CreateTable(tableName string, columns []column.Column, clusterIndex index.Index) error {
	if _, ok := db.table[tableName]; ok {
		return fmt.Errorf("this table [%s] is exist", tableName)
	}

	table, err := index.NewTable(columns, clusterIndex)
	db.table[tableName] = table

	fmt.Printf("New table [%s] is created\n", tableName)
	return err
}

func (db *Database) SelectTable(name string) *index.Table {
	return db.table[name]
}
