package linker

import (
	"root/index"
)

type Linker struct {
	link chan Payload
}

type Payload struct {
	data       []byte
	cluster    index.Cluster
	nonCluster []index.NonCluster
}

func New() Linker {
	return Linker{
		link: make(chan Payload),
	}
}

func (l *Linker) Send(data []byte, cluster index.Cluster, nonCluster []index.NonCluster) {
	l.link <- Payload{
		data:       data,
		cluster:    cluster,
		nonCluster: nonCluster,
	}
}

func (l *Linker) Receiver() ([]byte, index.Cluster, []index.NonCluster) {
	d := <-l.link
	return d.data, d.cluster, d.nonCluster
}
