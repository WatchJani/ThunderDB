package thunder

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"root/column"
	"root/cutter"
	"root/database"
	f "root/filter"
	"root/index"
	"root/linker"
	"root/table"
	"strings"
)

type Thunder struct {
	linker.Linker
	filePath string
	database map[string]*database.Database
}

func New() (*Thunder, error) {
	linker := linker.New()

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := path.Join(dir, "store.bin")
	cutter, err := cutter.New(linker, path, 10)
	if err != nil {
		log.Println(err)
	}

	go cutter.Cut()

	return &Thunder{
		Linker:   linker,
		filePath: path,
		database: make(map[string]*database.Database),
	}, nil
}

func (t *Thunder) NewDatabase(name string) error {
	t.database[name] = database.New(t.Linker)
	return nil
}

func (t *Thunder) NewTable(databaseName, tableName string, columns []column.Column) error {
	database := t.database[databaseName]
	return database.CreateTable(tableName, t.filePath, columns)
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

// ! provjeriti da li postoji kolona uopste ta koja se pretrazuje
func (t *Thunder) Search(databaseName, tableName string, data [][]byte) ([]byte, error) {
	database := t.database[databaseName]
	tableProcess := database.GetTable(tableName)

	dataSize := len(data)
	if dataSize%3 != 0 {
		return data[0], fmt.Errorf("wrong input")
	}

	filterField := make([]f.FilterField, dataSize/3)
	for columnIndex, dataIndex := 0, 0; dataIndex < len(data); dataIndex, columnIndex = dataIndex+3, columnIndex+1 {
		filterField[columnIndex] = f.New(data[dataIndex], data[dataIndex+1], data[dataIndex+2])
	}

	index, key := ChooseIndex(tableProcess, filterField)
	return index.Search(key, filterField, tableProcess.GetColumns())
}

func ChooseIndex(t *table.Table, filterField []f.FilterField) (index.Index, [][]byte) {
	for i, column := range filterField {
		userColumn := column.GetField()
		if userColumn == "id" {
			// var f int
			// if filterField[i].GetOperation() == "==" {
			// 	f = 1
			// }

			return t.GetClusterIndex(), [][]byte{filterField[i].GetValue()}
		}

		for j, index := range t.GetNonClusterIndex() {
			if index.GetByColumn()[0] == userColumn {
				key := ColumnBySearch(index.GetByColumn(), filterField) //crate key for index
				//f := 0
				// if filterField[len(key)-1].GetOperation() == "==" {
				// 	f = len(key)
				// }

				return t.GetNonClusterIndex()[j], key
			}
		}
	}

	return t.GetClusterIndex(), [][]byte{} //Work
}

func ColumnBySearch(index []string, filter []f.FilterField) [][]byte {
	key := make([][]byte, 0, len(index))
	for i := 0; i < len(index); i++ {
		found := false

		for j := i; j < len(filter); j++ {
			if index[i] == filter[j].GetField() {
				key = append(key, filter[j].GetValue())
				filter[j], filter[i] = filter[i], filter[j]

				if filter[j].GetOperation() != "==" {
					goto end
				}

				found = true
				break
			}
		}

		if !found {
			break
		}
	}

end:
	return key
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
	case "SEARCH":
		return t.CreateSearch(args)
	default:
		return nil, errors.New("command is not exist")
	}
}

func (t *Thunder) CreateDatabase(args []byte) ([]byte, error) {
	return nil, t.NewDatabase(string(args))
}

func (t *Thunder) CreateSearch(args []byte) ([]byte, error) {
	token := bytes.Split(args, []byte{' '})
	return t.Search(string(token[0]), string(token[1]), token[2:])
}

func (t *Thunder) CreateTable(args []byte) ([]byte, error) {
	token := strings.Split(string(args), " ")

	database, table := database.ParseDatabaseTable(token[0])
	t.NewTable(database, table, column.CreateColumn(token[1:]))

	return nil, nil
}

func (t *Thunder) CreateInsert(args []byte) ([]byte, error) {
	token := make([][]byte, 0, 3)

	start, pointer := 0, 0

	for end := 0; end < len(args); end++ {
		if args[end] == ' ' && pointer < 2 {
			token = append(token, args[start:end])
			start = end + 1
			pointer++
		}
	}

	if start < len(args) {
		token = append(token, args[start:])
	}

	if len(token) < 3 {
		return args, fmt.Errorf("wrong query input")
	}

	t.InsetData(string(token[0]), string(token[1]), token[2])

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
