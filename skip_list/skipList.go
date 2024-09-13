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
	down  *Node
	value int
	key   [][]byte
	time  time.Time
	leaf  bool
}

func NewNode(next, down *Node, value int, key [][]byte, leaf bool) Node {
	return Node{
		time:  time.Now(),
		next:  next,
		down:  down,
		value: value,
		key:   key,
		leaf:  leaf,
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

	index := 0
	// compareFn := GetCompareFuncType(s.dataType[index])
	for {
		for current.next != nil {
			// num := compareFn(current.next.key[index], key[index])
			num := bytes.Compare(current.next.key[index], key[index])
			if num == -1 {
				current = current.next
			} else if num == 0 {
				if index+1 < len(key) {
					index++
					// compareFn = GetCompareFuncType(s.dataType[index])
				} else {
					break // i found the key
				}
			} else {
				break
			}
		}

		if current.leaf {
			break
		}

		stack.Push(current)
		current = current.down
	}

	nextNode := current.next

	node := s.Pool.Insert()

	current.next = node
	*node = NewNode(nextNode, nil, value, key, true) // create new leaf node

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
		*node = NewNode(nextNode, downNode, value, key, false) // create new internal node
	}

	stack.Clear()       //Clear stack
	s.Stack.Push(stack) // return to stack stack
}

func flipCoin(percentage float64) bool {
	return rand.Float64() < percentage
}

func (s *SkipList) Search(key [][]byte) (bool, int) {
	s.Lock()
	defer s.Unlock()

	current := s.roots[s.rootIndex]

	index := 0
	// compareFn := GetCompareFuncType(s.dataType[index])

	for {
		for current.next != nil {
			// num := compareFn(current.next.key[index], key[index])
			num := bytes.Compare(current.next.key[index], key[index])
			if num == -1 { // < n
				current = current.next
			} else if num == 0 { // == n
				if index+1 < len(key) {
					index++
					// compareFn = GetCompareFuncType(s.dataType[index])
				} else { // > n
					return true, current.next.value
				}
			} else {
				break
			}
		}

		if current.leaf {
			return false, current.next.value
		}

		current = current.down
	}
}

func (s *SkipList) Read() {
	for startNode := s.roots[0].next; startNode != nil; startNode = startNode.next {
		fmt.Println(startNode)
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

// func CompareFloat(key, currentKey []byte) int {
// 	keyFromString, _ := strconv.ParseFloat(string(key), 64)
// 	currentKeyFromString, _ := strconv.ParseFloat(string(currentKey), 64)

// 	return cmp.Compare(keyFromString, currentKeyFromString)
// }

// func CompareInt(key, currentKey []byte) int {
// 	keyFromString, _ := strconv.Atoi(string(currentKey))
// 	currentKeyFromString, _ := strconv.Atoi(string(currentKey))

// 	return cmp.Compare(keyFromString, currentKeyFromString)
// }

// func CompareOtherType(key, currentKey []byte) int {
// 	return bytes.Compare(key, currentKey)
// }

// func GetCompareFuncType(dataType string) func([]byte, []byte) int {
// 	switch dataType {
// 	case "INT":
// 		return CompareInt
// 	case "FLOAT":
// 		return CompareFloat
// 	default:
// 		return CompareOtherType
// 	}
// }
