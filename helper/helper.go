package helper

import (
	"root/column"
	"strconv"
)

func ReadSingleData(data []byte, columns []column.Column) ([][]byte, error) {
	columnData := make([][]byte, len(columns))

	index := 0
	for i := range columns {
		index += 5
		size := data[index-5 : index]

		num, err := strconv.Atoi(string(size))
		if err != nil {
			return columnData, err
		}

		end := index + num
		columnData[i] = data[index:end]
		index += num
	}

	return columnData, nil
}

func GetColumnNameIndex(name string, columns []column.Column) int {
	for index, column := range columns {
		if column.GetName() == name {
			return index
		}
	}

	return -1
}
