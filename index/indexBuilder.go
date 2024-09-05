package index

import (
	"fmt"
)

const Cluster string = "cluster"

type IndexBuilder struct {
	columns []string
	index   []Index
}

func (ib *IndexBuilder) AddIndex(newIndex Index) error {
	for _, index := range ib.columns {
		if index == newIndex.byColumn[0] {
			ib.index = append(ib.index, newIndex)
			return nil
		}
	}

	return fmt.Errorf("column %s not exist", newIndex.name)
}

func NewIndexBuilder(columns []string, index []Index) (*IndexBuilder, error) {
	// counter, clusterIndex := 0, index[0]
	// for i := 0; i < len(index); i++ {
	// 	for _, column := range columns {
	// 		if clusterIndex.byColumn[i] == column {
	// 			counter++
	// 			break
	// 		}
	// 	}

	// 	if i != counter {
	// 		return nil, fmt.Errorf("index column field %s not exit", index[0].byColumn[i])
	// 	}
	// }

	return &IndexBuilder{
		columns: columns,
		index:   index,
	}, nil
}

func (ib *IndexBuilder) Choice(userField []string) (Index, bool) {
	var i, j int

	//check cluster index
	clusterIndex := ib.index[0]
	for i < clusterIndex.GetColumnNumber() && j < len(userField) {
		for j < len(userField) {
			if clusterIndex.byColumn[i] == userField[j] {
				userField[i], userField[j] = userField[j], userField[i]
				i++
				j = i
				break
			}
			j++
		}

		if i != j {
			break
		}
	}

	if i != 0 {
		return clusterIndex, false
	}

	//check non-cluster index
	for i := 1; i < len(ib.index); i++ {
		for j := 0; j < len(userField); j++ {
			if ib.index[i].name == userField[j] {
				return ib.index[i], false
			}
		}
	}

	return clusterIndex, true
}
