package database

import "root/index"

type Database struct {
	table map[string]*index.IndexBuilder
}
