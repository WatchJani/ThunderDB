package linker

import (
	"root/column"
	"root/index"
	"sync"
)

type Linker struct {
	Link chan Payload
}

func New() Linker {
	return Linker{
		Link: make(chan Payload),
	}
}

type Payload struct {
	memTable   *[]byte
	end        int
	nonCluster []*index.NonCluster
	cluster    *index.Cluster
	columns    []column.Column
	*sync.WaitGroup
}

func (l *Linker) Receive() (*[]byte, int, *sync.WaitGroup, []*index.NonCluster, *index.Cluster, []column.Column) {
	payload := <-l.Link
	return payload.memTable, payload.end, payload.WaitGroup, payload.nonCluster, payload.cluster, payload.columns
}

func NewPayload(memTable *[]byte, end int, wg *sync.WaitGroup, nonCluster []*index.NonCluster, cluster *index.Cluster, columns []column.Column) Payload {
	return Payload{
		memTable:   memTable,
		end:        end,
		WaitGroup:  wg,
		nonCluster: nonCluster,
		columns:    columns,
		cluster:    cluster,
	}
}
