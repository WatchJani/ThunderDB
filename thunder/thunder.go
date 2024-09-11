package thunder

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"root/column"
	"root/database"
	"root/index"
	"root/table"
	"strconv"
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

func (t *Thunder) CreateDatabase(name []byte) error {
	t.Lock()
	defer t.Unlock()

	databaseName := string(name)

	if _, ok := t.thunder[databaseName]; ok {
		return fmt.Errorf("database [%s] is exist", name)
	}

	t.thunder[databaseName] = database.New(databaseName)
	fmt.Printf("New database [%s] is created\n", name)
	return nil
}

func (t *Thunder) Insert(query []byte) error {
	parts := bytes.SplitN(query, []byte(" "), 2)
	if len(parts) < 2 {
		return fmt.Errorf("wrong query: %s", query)
	}

	dbTable := bytes.Split(parts[0], []byte("."))
	if len(dbTable) != 2 {
		return fmt.Errorf("wrong database name or table name: %s", parts[0])
	}

	database := string(dbTable[0])
	queryDatabase, err := t.SelectDatabase(database)
	if err != nil {
		return err
	}

	tableReq := string(dbTable[1])
	queryTable, err := queryDatabase.SelectTable(tableReq)
	if err != nil {
		return err
	}

	conditionsPart := parts[1]

	// fmt.Println(queryTable.GetColumn())
	// fmt.Println(queryTable.GetIndex())
	// fmt.Println(string(conditionsPart))

	//copy data to memTable
	memTableOffset := queryTable.Write(conditionsPart)
	offset := make([]int, 0, queryTable.GetColumnNum()*2)

	index := 0
	for range queryTable.GetColumn() {
		index += 5
		size := conditionsPart[index-5 : index]

		num, err := strconv.Atoi(string(size))
		if err != nil {
			return err
		}

		start := memTableOffset + index
		end := start + num

		offset = append(offset, start, end)
		// fmt.Println(column, string(conditionsPart[index:index+num]))
		index += num
	}

	// fmt.Println(offset)
	fmt.Println(string(queryTable.GetData()))

	//!ima bug negdje!!!
	for _, index := range queryTable.GetIndex() {

		indexColumnNumber := index.GetColumnNumber()
		key := make([][]byte, 0, indexColumnNumber)
		
		for j := 0; j < indexColumnNumber; j++ {

			position, err := queryTable.FindIndexColumn(index.GetByColumn(j))
			if err != nil {
				return err
			}

			key = append(key, queryTable.GetDataOnSpecificPosition(
				offset[position],
				offset[position+1],
			))
		}

		index.Insert(key, memTableOffset)
	}

	return nil
}

func (t *Thunder) Search(query []byte) error {
	parts := bytes.SplitN(query, []byte(" "), 2)
	if len(parts) < 2 {
		return fmt.Errorf("wrong query: %s", query)
	}

	dbTable := bytes.Split(parts[0], []byte("."))
	if len(dbTable) != 2 {
		return fmt.Errorf("wrong database name or table name: %s", parts[0])
	}

	database := string(dbTable[0])
	queryDatabase, err := t.SelectDatabase(database)
	if err != nil {
		return err
	}

	tableReq := string(dbTable[1])
	queryTable, err := queryDatabase.SelectTable(tableReq)
	if err != nil {
		return err
	}

	conditionsPart := parts[1]
	conditionStrings := bytes.Split(conditionsPart, []byte(" "))
	conditions := make([]table.Condition, 0, 2)

	for i := 0; i < len(conditionStrings); i += 4 {
		Field := conditionStrings[i]
		Type := conditionStrings[i+1]
		Operator := conditionStrings[i+2]
		Value := conditionStrings[i+3]

		conditions = append(conditions, table.Condition{
			Field:    string(Field),
			Type:     string(Type),
			Operator: Operator,
			Value:    Value,
		})
	}

	return queryTable.Search(conditions)
}

func (t *Thunder) QueryParser(payload []byte) error {
	command, args := findCommand(payload)
	switch command {
	case "CREATE_DATABASE":
		return t.CreateDatabase(args)
	case "CREATE_TABLE":
		return t.CreateTable(args)
	case "INSERT":
		return t.Insert(args)
	case "SEARCH":
		return t.Search(args)
	default:
		return errors.New("command is not exist")
	}
}

func (t *Thunder) CreateTable(tableQuery []byte) error {
	query := string(tableQuery)
	tableInfo := table.NewTableInfo()

	re := regexp.MustCompile(`(?i)(\w+)\.(\w+)\s+(.*)\s*\[(.*?)\]`)

	matches := re.FindStringSubmatch(query)
	if len(matches) == 0 {
		return fmt.Errorf("wrong SQL query format")
	}

	database := t.thunder[matches[1]]
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

	//! set default cluster index
	indexColumns := strings.Split(matches[4], ",")
	for _, index := range indexColumns {
		tableInfo.Indexes = append(tableInfo.Indexes, strings.TrimSpace(index))
	}

	return database.CreateTable(tableInfo.Table, tableInfo.Columns, index.NewClusterIndex(tableInfo.Indexes...))
}

func findCommand(payload []byte) (string, []byte) {
	for i := 0; i < len(payload); i++ {
		if payload[i] == ' ' {
			return string(payload[:i]), payload[i+1:]
		}
	}

	return "", []byte{}
}

func (t *Thunder) SelectDatabase(name string) (*database.Database, error) {
	if database, ok := t.thunder[name]; ok {
		return database, nil
	}

	return nil, fmt.Errorf("this database [%s] is not exist", name)
}
