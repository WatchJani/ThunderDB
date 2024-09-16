package main

import (
	"log"
	"root/query"
	"root/thunder"
	"time"
)

// type Node struct {
// 	buff chan []byte
// }

// func (n *Node) Send(data []byte) {
// 	n.buff <- data
// }

// func (n *Node) Receive() []byte {
// 	return <-n.buff
// }

// func New() *Node {
// 	return &Node{
// 		buff: make(chan []byte),
// 	}
// }

// func (n *Node) Reader() {
// 	for {
// 		data := n.Receive()
// 		fmt.Println(data[:10])

// 		time.Sleep(5 * time.Second)
// 	}
// }

func main() {

	thunder := thunder.New()
	if err := thunder.NewDatabase("netflix"); err != nil {
		log.Println(err)
	}

	if err := thunder.NewTable("netflix", "user"); err != nil {
		log.Println(err)
	}

	for range 171197 * 2 {
		thunder.Inset("netflix", "user", query.Insert()[20:])
	}

	time.Sleep(5 * time.Second)

}
