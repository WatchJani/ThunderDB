package linker

import (
	"root/column"
	"root/index"
)

type Linker struct {
	link chan Payload
}

type Payload struct {
	data       []byte
	cluster    index.Cluster
	nonCluster []index.NonCluster
	columns    []column.Column
}

func New() Linker {
	return Linker{
		link: make(chan Payload),
	}
}

func (l *Linker) Send(data []byte, cluster index.Cluster, nonCluster []index.NonCluster, columns []column.Column) {
	l.link <- Payload{
		data:       data,
		cluster:    cluster,
		nonCluster: nonCluster,
		columns:    columns,
	}
}

func (l *Linker) Receiver() ([]byte, index.Cluster, []index.NonCluster, []column.Column) {
	d := <-l.link
	return d.data, d.cluster, d.nonCluster, d.columns
}
