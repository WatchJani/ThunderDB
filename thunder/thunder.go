package thunder

import (
	"errors"
	"fmt"
	"regexp"
	"root/column"
	"root/database"
	"root/index"
	"strings"
	"sync"
)

type Thunder struct {
	thunder map[string]*database.Database
	sync.RWMutex
}

func New() *Thunder {
	return &Thunder{
		thunder: make(map[string]*database.Database),
	}
}

func (t *Thunder) CreateDatabase(name string) error {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.thunder[name]; ok {
		return fmt.Errorf("database [%s] is exist", name)
	}

	t.thunder[name] = database.New(name)
	fmt.Printf("New database [%s] is created\n", name)
	return nil
}

// func (t *Thunder) Write() {

// }

// func (t *Thunder) Read() {

// }

func (t *Thunder) QueryParser(payload []byte) error {
	command, args := findCommand(payload)

	switch command {
	case "CREATE_DATABASE":
		return t.CreateDatabase(args)
	case "CREATE_TABLE":
		return t.CreateTable(args)
	case "INSERT":
		return nil
	default:
		return errors.New("command is not exist")
	}
}

func (t *Thunder) CreateTable(query string) error {
	tableInfo := index.NewTableInfo()

	re := regexp.MustCompile(`(?i)(\w+)\.(\w+)\s+(.*)\s*\[(.*?)\]`)

	matches := re.FindStringSubmatch(query)
	if len(matches) == 0 {
		return fmt.Errorf("wrong SQL query format")
	}

	tableInfo.Database = matches[1]
	tableInfo.Table = matches[2]

	columnsPart := matches[3]
	columnDefs := strings.Split(columnsPart, ",")
	for _, colDef := range columnDefs {
		colDef = strings.TrimSpace(colDef)
		colParts := strings.Split(colDef, " ")
		if len(colParts) != 2 {
			return fmt.Errorf("wrong column: %s", colDef)
		}

		tableInfo.Columns = append(tableInfo.Columns, column.New(colParts[0], colParts[1]))
	}

	indexColumns := strings.Split(matches[4], ",")
	for _, index := range indexColumns {
		tableInfo.Indexes = append(tableInfo.Indexes, strings.TrimSpace(index))
	}

	fmt.Println(tableInfo.String())

	return nil
}

func findCommand(payload []byte) (string, string) {
	for i := 0; i < len(payload); i++ {
		if payload[i] == ' ' {
			return string(payload[:i]), string(payload[i+1:])
		}
	}

	return "", ""
}
