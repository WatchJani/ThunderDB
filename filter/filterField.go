package filter

type FilterField struct {
	filter    Filter
	field     string
	operation string
	value     []byte
}

func New(field, operation, input []byte) FilterField {
	return FilterField{
		filter: GenerateFilter(
			operation,
			input,
		),
		field:     string(field),
		operation: string(operation),
		value:     input,
	}
}

func (f *FilterField) GetField() string {
	return f.field
}

func (f *FilterField) GetOperation() string {
	return f.operation
}

func (f *FilterField) GetValue() []byte {
	return f.value
}
