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
