package cutter

import "root/linker"

type Cutter struct {
	linker.Linker
}

func New(linker linker.Linker) Cutter {
	return Cutter{
		Linker: linker,
	}
}

func (c *Cutter) Cut() {
	for {
		c.Receiver()
	}
}
