package thunder

import (
	"errors"
	"fmt"
	"log"
	"root/column"
	"root/cutter"
	"root/database"
	"root/linker"
	"strings"
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

func (t *Thunder) NewTable(databaseName, tableName string, columns []column.Column) error {
	database := t.database[databaseName]

	err := database.CreateTable(tableName, columns)
	fmt.Printf("Table %s is created\n", tableName)

	return err
}

func (t *Thunder) InsetData(databaseName, tableName string, data []byte) {
	database := t.database[databaseName]
	table := database.GetTable(tableName)

	memTableOffset := table.Insert(data) //write data to memTable or send on disk

	//update nonClusterIndex

	_ = memTableOffset
}

func (t *Thunder) QueryParser(payload []byte) ([]byte, error) {
	command, args := findCommand(payload)
	switch command {
	case "CREATE_DATABASE":
		return t.CreateDatabase(args)
	case "CREATE_TABLE":
		return t.CreateTable(args)
	case "INSERT":
		return t.CreateInsert(args)
	// case "SEARCH":
	// return t.Search(args)
	default:
		return nil, errors.New("command is not exist")
	}
}

func (t *Thunder) CreateDatabase(args []byte) ([]byte, error) {
	err := t.NewDatabase(string(args))
	fmt.Printf("Database %s is created\n", args)

	return nil, err
}

func (t *Thunder) CreateTable(args []byte) ([]byte, error) {
	token := strings.Split(string(args), " ")

	database, table := database.ParseDatabaseTable(token[0])
	t.NewTable(database, table, column.CreateColumn(token[1:]))

	return nil, nil
}

// not secure
func (t *Thunder) CreateInsert(args []byte) ([]byte, error) {
	si, store := []byte{'.', ' '}, make([]string, 2)

	for index, siCounter := 0, 0; index < len(args); index++ {
		if args[index] == si[siCounter] {
			store[siCounter] = string(args[:index])
			args, index = args[index:], 0
			siCounter++

			if args[index] == ' ' {
				break
			}
		}
	}

	t.InsetData(store[0], store[1][1:], args[1:])

	return nil, nil
}

func findCommand(payload []byte) (string, []byte) {
	for i := 0; i < len(payload); i++ {
		if payload[i] == ' ' {
			return string(payload[:i]), payload[i+1:]
		}
	}

	return "", []byte{}
}
