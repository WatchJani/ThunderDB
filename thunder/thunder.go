package thunder

import (
	"errors"
	"log"
	"root/column"
	"root/cutter"
	"root/database"
	"root/linker"
	"root/table"
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
	return database.CreateTable(tableName, columns)
}

func (t *Thunder) NewIndex(databaseName, tableName string, columns []string) error {
	database := t.database[databaseName]
	table := database.GetTable(tableName)

	return table.NewIndex(columns...)
}

func (t *Thunder) InsetData(databaseName, tableName string, data []byte) {
	database := t.database[databaseName]
	tableProcess := database.GetTable(tableName)

	memTableOffset := tableProcess.Insert(data) //write data to memTable or send on disk

	columnData, err := tableProcess.ReadSingleData(data[5:])
	if err != nil {
		log.Println(err)
	}

	for _, index := range tableProcess.GetIndexes() {
		key := table.GenerateKey(index, columnData, tableProcess.GetColumns())
		index.Insert(key, memTableOffset)
	}
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
	case "INDEX":
		return t.CreateIndex(args)
	// case "SEARCH":
	// return t.Search(args)
	default:
		return nil, errors.New("command is not exist")
	}
}

func (t *Thunder) CreateDatabase(args []byte) ([]byte, error) {
	return nil, t.NewDatabase(string(args))
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
			store[siCounter], args, index, siCounter = string(args[:index]), args[index:], 0, siCounter+1

			if args[index] == ' ' {
				break
			}
		}
	}

	t.InsetData(store[0], store[1][1:], args[1:])

	return nil, nil
}

func (t *Thunder) CreateIndex(args []byte) ([]byte, error) {
	token := strings.Split(string(args), " ")
	database, table := database.ParseDatabaseTable(token[0])

	return nil, t.NewIndex(database, table, token[1:])
}

func findCommand(payload []byte) (string, []byte) {
	for i := 0; i < len(payload); i++ {
		if payload[i] == ' ' {
			return string(payload[:i]), payload[i+1:]
		}
	}

	return "", []byte{}
}
