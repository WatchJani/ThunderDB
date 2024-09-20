package skip_list

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/WatchJani/pool"
	st "github.com/WatchJani/stack"
)

type SkipList struct {
	roots     []*Node
	rootIndex int
	st.Stack[st.Stack[*Node]]
	pool.Pool[Node]
	sync.RWMutex
	percentage float64
	height     int

	//If dataTYpe is int or float, then data will be searched different
	// dataType []string
}

func New(height, capacity int, percentage float64) *SkipList {
	//fix this part to be dynamic
	stack := st.New[st.Stack[*Node]](250) // max number of parallel readings 250

	for range 250 {
		stack.Push(st.New[*Node](height))
	}

	roots := make([]*Node, height)
	prevues := &Node{leaf: true}
	roots[0] = prevues

	for index := 1; index < len(roots); index++ {
		roots[index] = &Node{down: prevues}
		prevues = roots[index]
	}

	return &SkipList{
		roots:      roots,
		Stack:      stack,
		Pool:       pool.New[Node](capacity),
		percentage: percentage,
		height:     height,
		// dataType:   dataType,
	}
}

type Node struct {
	next  *Node
	back  *Node
	down  *Node
	value int
	key   [][]byte
	time  time.Time
	leaf  bool
}

func NewNode(back, next, down *Node, value int, key [][]byte, leaf bool) Node {
	return Node{
		time:  time.Now(),
		next:  next,
		down:  down,
		value: value,
		key:   key,
		leaf:  leaf,
		back:  back,
	}
}

func (s *SkipList) Insert(key [][]byte, value int) {
	s.Lock()
	defer s.Unlock()

	current := s.roots[s.rootIndex]

	stack, err := s.Stack.Pop()
	if err != nil {
		stack = st.New[*Node](s.height)
	}

	for { //Down
		current = HorizontalSearch(current, key)

		if current.leaf {
			break
		}

		stack.Push(current)
		current = current.down
	}

	nextNode := current.next

	node := s.Pool.Insert()

	current.next = node

	if nextNode != nil {
		nextNode.back = node
	}

	*node = NewNode(current, nextNode, nil, value, key, true) // create new leaf node

	for flipCoin(s.percentage) {
		downNode := node
		leftNode, err := stack.Pop()

		if err != nil {
			if s.rootIndex+1 > s.height {
				break
			}

			s.rootIndex++
			leftNode = s.roots[s.rootIndex]
		}

		nextNode = leftNode.next

		node = s.Pool.Insert()
		leftNode.next = node
		*node = NewNode(leftNode, nextNode, downNode, value, key, false) // create new internal node
	}

	stack.Clear()       //Clear stack
	s.Stack.Push(stack) // return to stack stack
}

func flipCoin(percentage float64) bool {
	return rand.Float64() < percentage
}

func HorizontalSearch(current *Node, key [][]byte) *Node {
	for current.next != nil {
		var num int
		for i := 0; i < len(key); i++ { //Search for key
			num = bytes.Compare(current.next.key[i], key[i])

			if num != 0 {
				break
			}
		}

		if num == -1 {
			current = current.next
		} else {
			break
		}
	}

	return current
}

func (s *SkipList) Search(key [][]byte, operation string) (*Node, bool) {
	s.Lock()
	defer s.Unlock()

	current := s.roots[s.rootIndex]

	for {
		current = HorizontalSearch(current, key)

		if current.leaf {
			switch operation {
			case "<":
				if current == s.roots[0] {
					return nil, false
				}
				return current, false
			default:
				return current.next, false
			}
		}

		current = current.down
	}
}

func (s *SkipList) RootNode() *Node {
	return s.roots[0].next
}

func (s *Node) NextNode() *Node {
	return s.next
}

func (n *Node) Key() [][]byte {
	return n.key
}

func (n *Node) GetValue() int {
	return n.value
}

func (s *SkipList) Clear() {
	for index := range s.roots {
		s.roots[index].next = nil
	}

	s.Pool.Clear()
}

// for test
func (s *SkipList) ReadAllFromLeftToRight() {
	for root := s.roots[0].next; root != nil; root = root.next {
		fmt.Println(root)
	}
}

// for test
func (s *SkipList) ReadAllFromRightToLeft() {
	root := s.roots[0].next

	//go to last right position
	for root.next != nil {
		root = root.next
	}

	for root.back != nil {
		fmt.Println(root)
		root = root.back
	}
}
