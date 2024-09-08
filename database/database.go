package database

import (
	"fmt"
	"root/column"
	"root/index"
	"root/table"
)

type Database struct {
	table map[string]*table.Table //index builder is table
}

func New(name string) *Database {
	return &Database{
		table: make(map[string]*table.Table),
	}
}

// name of table, and all information about table
func (db *Database) CreateTable(tableName string, columns []column.Column, clusterIndex index.Index) error {
	if _, ok := db.table[tableName]; ok {
		return fmt.Errorf("this table [%s] is exist", tableName)
	}

	table, err := table.NewTable(columns, clusterIndex)
	db.table[tableName] = table

	fmt.Printf("New table [%s] is created\n", tableName)
	return err
}

func (db *Database) SelectTable(tableName string) (*table.Table, error) {
	if table, ok := db.table[tableName]; ok {
		return table, nil
	}

	return nil, fmt.Errorf("this table [%s] is not exist", tableName)
}
