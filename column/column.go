package column

type Column struct {
	name     string
	dataType string
}

func New(name, dataType string) Column {
	return Column{
		name:     name,
		dataType: dataType,
	}
}

func (c *Column) GetName() string {
	return c.name
}

func (c *Column) GetDataType() string {
	return c.dataType
}

func DefaultType() Column {
	return Column{
		name:     "id",
		dataType: "UUID",
	}
}

func CreateColumn(token []string) []Column {
	columns := make([]Column, len(token)/2+1)
	columns[0] = DefaultType()

	for index, counter := 0, 1; index < len(token); index += 2 {
		columns[counter] = New(token[index], token[index+1])
		counter++
	}

	return columns
}
