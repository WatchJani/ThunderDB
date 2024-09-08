package BPTree

import (
	"errors"
	"fmt"
)

type KeyType interface {
	int | string | float64 | float32 | int16 | int8 | int32 | int64
}

type Tree[K KeyType, V any] struct {
	root   *Node[K, V]
	degree int
}

type item[K KeyType, V any] struct {
	key   K
	value V
}

type Node[K KeyType, V any] struct {
	items     []item[K, V]
	Children  []*Node[K, V]
	nextNodeL *Node[K, V]
	nextNodeR *Node[K, V]
	pointer   int
}

type Stack[K KeyType, V any] struct {
	store []positionStr[K, V]
}

type positionStr[K KeyType, V any] struct {
	node     *Node[K, V]
	position int
}

func newItem[K KeyType, V any](key K, value V) item[K, V] {
	return item[K, V]{
		key:   key,
		value: value,
	}
}

func newNode[K KeyType, V any](degree int) Node[K, V] {
	return Node[K, V]{
		items:    make([]item[K, V], degree+1),
		Children: make([]*Node[K, V], degree+2),
	}
}

func (n *Node[K, V]) delete() {
	n.items = nil
	n.Children = nil
	n.pointer = 0
	n.nextNodeL = nil
	n.nextNodeR = nil
}

func New[K KeyType, V any](degree int) *Tree[K, V] {
	if degree < 3 {
		degree = 3
	}

	return &Tree[K, V]{
		degree: degree,
		root: &Node[K, V]{
			items:    make([]item[K, V], degree+1),
			Children: make([]*Node[K, V], degree+2),
		},
	}
}

func (n *Node[K, V]) search(target K) (int, bool) {
	low, high := 0, n.pointer-1

	for low <= high {
		mid := (low + high) / 2

		if n.items[mid].key == target {
			return mid + 1, true
		} else if n.items[mid].key < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low, false
}

func newStack[K KeyType, V any]() Stack[K, V] {
	return Stack[K, V]{
		store: make([]positionStr[K, V], 0, 4),
	}
}

func (s *Stack[K, V]) Push(node *Node[K, V], position int) {
	s.store = append(s.store, positionStr[K, V]{
		node:     node,
		position: position,
	})
}

func (s *Stack[K, V]) Pop() (positionStr[K, V], error) {
	if len(s.store) == 0 {
		return positionStr[K, V]{}, errors.New("stack is empty")
	}

	pop := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return pop, nil
}

func (t *Tree[K, V]) Find(key K) (V, error) {
	for next := t.root; next != nil; {
		index, found := next.search(key)

		if found {
			return next.items[index-1].value, nil
		}

		next = next.Children[index]
	}

	var res V
	return res, fmt.Errorf("key %v not found", key)
}

func (t *Tree[K, V]) Insert(key K, value V) {
	stack, item := newStack[K, V](), newItem(key, value)
	position, found := findLeaf(t.root, &stack, key)

	current, _ := stack.Pop()
	//update just state state
	if found {
		current.node.items[position].value = value
		return
	}

	if middleKey, nodeChildren := insertLeaf(current.node, position, t.degree, item); nodeChildren != nil {
		for {
			temp := current

			stack, err := stack.Pop()
			if err != nil {
				current = temp
				break
			}

			parent := stack.node

			parent.pointer += insert(parent.items, middleKey, stack.position)
			chIndex := childrenIndex(middleKey.key, parent.items[stack.position].key, stack.position)

			insert(parent.Children, nodeChildren, chIndex)
			if parent.pointer < t.degree {
				return
			}

			middle := parent.pointer / 2
			middleKey = parent.items[middle]

			newNode := newNode[K, V](t.degree)
			newNode.pointer += migrate(newNode.items, parent.items[:middle], 0) //migrate half element to left child node
			migrate(newNode.Children, parent.Children[:middle+1], 0)

			parent.pointer -= deleteElement(parent.items, 0, parent.pointer-middle+1-t.degree&1)
			migrate(parent.Children, parent.Children[middle+1:], 0)
			nodeChildren = &newNode

			current = stack
		}

		rootNode := newNode[K, V](t.degree)
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.Children[0] = nodeChildren
		rootNode.Children[1] = current.node
		t.root = &rootNode
	}
}

func childrenIndex[K KeyType](key, value K, index int) int {
	if value < key {
		return index + 1
	}

	return index
}

func insert[T any](list []T, insert T, position int) int {
	copy(list[position+1:], list[position:])
	return copy(list[position:], []T{insert})
}

func migrate[T any](list, migrateElement []T, position int) int {
	return copy(list[position:], migrateElement)
}

func deleteElement[T any](list []T, position, deletion int) int {
	copy(list[position:], list[position+deletion:])
	return deletion
}

func insertLeaf[K KeyType, V any](current *Node[K, V], position, degree int, item item[K, V]) (item[K, V], *Node[K, V]) {
	current.pointer += insert(current.items, item, position)

	if current.pointer < degree {
		return item, nil
	}

	newNode := newNode[K, V](degree)
	middle := degree / 2

	newNode.pointer += migrate(newNode.items, current.items[:middle], 0)
	current.pointer -= deleteElement(current.items, 0, current.pointer-middle-degree%2)

	//update links between leafs
	if current.nextNodeL != nil {
		newNode.nextNodeL = current.nextNodeL
		current.nextNodeL.nextNodeR = &newNode
	}

	current.nextNodeL = &newNode
	newNode.nextNodeR = current

	return current.items[0], &newNode
}

func minAllowed(degree, numElement int) bool {
	return (degree-1)/2 <= numElement
}

func findLeaf[K KeyType, V any](root *Node[K, V], stack *Stack[K, V], key K) (int, bool) {
	position, found := 0, false

	for current := root; current != nil; {
		position, found = current.search(key)
		stack.Push(current, position)

		current = current.Children[position]
	}

	return position, found
}

func (t *Tree[K, V]) Delete(key K) error {
	stack := newStack[K, V]()
	_, found := findLeaf(t.root, &stack, key)

	if !found {
		return fmt.Errorf("key %v is not exist", key)
	}

	current, _ := stack.Pop()

	for {
		current.node.pointer -= deleteElement(current.node.items, indexElement(current.position), 1)

		if minAllowed(t.degree, current.node.pointer) || (found && t.root.pointer == 0) {
			return nil
		}

		temp := current //current
		parent, err := stack.Pop()
		if err != nil {
			break
		}

		sibling, side, operation := sibling(parent, t.degree)

		if operation {
			transfer(parent, temp, sibling, found, side)
			return nil
		} else {
			merge(temp.node, sibling, parent, found, side)
		}

		if found {
			found = !found
		}

		current = parent
	}

	if t.root.pointer == 0 && len(t.root.Children) > 0 {
		if t.root.Children[0] != nil {
			t.root = t.root.Children[0]
		} else {
			t.root = t.root.Children[1]
		}
	}

	return nil
}

func siblingExist[K KeyType, V any](parent positionStr[K, V], index int) (*Node[K, V], bool) {
	index = parent.position + index

	if index < 0 || index > parent.node.pointer {
		return nil, false
	}

	sibling := parent.node.Children[index]
	return sibling, true
}

func sibling[K KeyType, V any](parent positionStr[K, V], degree int) (*Node[K, V], bool, bool) {
	var (
		potential *Node[K, V]
		side      bool
	)

	if sibling, isExist := siblingExist(parent, -1); isExist {
		if minAllowed(degree, sibling.pointer-1) {
			return sibling, true, true
		}

		potential, side = sibling, true
	}

	if sibling, isExist := siblingExist(parent, +1); isExist {
		if minAllowed(degree, sibling.pointer-1) {
			return sibling, false, true
		}

		if potential == nil {
			potential, side = sibling, false
		}
	}

	return potential, side, false
}

func sideFn(side bool, pointer int) int {
	if side {
		return 0
	}

	return pointer
}

func indexElement(index int) int {
	if index > 0 {
		return index - 1
	}

	return index
}

func merge[K KeyType, V any](current, sibling *Node[K, V], parent positionStr[K, V], leafInternal, side bool) {
	parentElement := parent.node.items[indexElement(parent.position)]
	position := sideFn(side, current.pointer)

	if leafInternal {
		if current.nextNodeL == sibling {
			current.nextNodeL = sibling.nextNodeL
			if sibling.nextNodeL != nil {
				sibling.nextNodeL.nextNodeR = current
			}
		} else {
			current.nextNodeR = sibling.nextNodeR
			if sibling.nextNodeR != nil {
				sibling.nextNodeR.nextNodeL = current
			}
		}
	} else {
		current.pointer += insert(current.items, parentElement, position)

		if !side {
			position++
		}

		insertSet(current.Children, sibling.Children[:sibling.pointer+1], position)
	}

	if side {
		deleteElement(parent.node.Children, parent.position-1, 1)
	} else {
		deleteElement(parent.node.Children, parent.position+1, 1)
	}

	current.pointer += insertSet(current.items, sibling.items[:sibling.pointer], position)

	sibling.delete()
}

func insertSet[T any](list []T, insert []T, position int) int {
	copy(list[position+len(insert):], list[position:])
	return copy(list[position:], insert)
}

func transfer[K KeyType, V any](parent, current positionStr[K, V], sibling *Node[K, V], leafInternal, side bool) {
	itemIndex := 0
	parentPosition := parent.position - 1
	childInsertPosition := current.node.pointer
	insertPosition := current.node.pointer

	if side {
		itemIndex = sibling.pointer - 1
		childInsertPosition = 0
		insertPosition = 0
	} else {
		parentPosition++
	}

	if leafInternal {
		siblingItem := sibling.items[itemIndex]
		if !side {
			siblingItem = sibling.items[1]
		}
		parent.node.items[parentPosition] = siblingItem
		current.node.pointer += insert(current.node.items, sibling.items[itemIndex], insertPosition)
	} else {
		current.node.pointer += insert(current.node.items, parent.node.items[parentPosition], childInsertPosition)
		parent.node.items[parentPosition] = sibling.items[itemIndex]

		if !side {
			insert(current.node.Children, sibling.Children[0], current.node.pointer) //check right side -> itemIndex+1(work for left) -> itemIndex
			deleteElement(sibling.Children, 0, 1)
		} else {
			insert(current.node.Children, sibling.Children[itemIndex+1], childInsertPosition) //check right side -> itemIndex+1(work for left) -> itemIndex
		}
	}

	sibling.pointer -= deleteElement(sibling.items, itemIndex, 1)
}

func (t *Tree[K, V]) TestFunc() int {
	current := t.root
	for current.Children[0] != nil {
		current = current.Children[0]
	}

	var counter int

	for current != nil {
		for _, value := range current.items[:current.pointer] {
			counter++
			// fmt.Println(counter, value)
			_ = value
		}
		// fmt.Println("======")
		current = current.nextNodeR
	}

	return counter
}

func (t *Tree[K, V]) GetRoot() *Node[K, V] {
	return t.root
}

func (t *Tree[K, V]) RangeUp(source, destination K, command string) []item[K, V] {
	var (
		index   int
		prevues *Node[K, V]
		src     []item[K, V] = make([]item[K, V], 0, 10)
	)

	for next := t.root; next != nil; {
		index, _ := next.search(source)
		prevues, next = next, next.Children[index]
	}

	index++

	fmt.Println(prevues.items[index])

	rangeFn, nextFn := CommandRange[K, V](command)

	for prevues != nil {
		for index < prevues.pointer {
			if rangeFn(index, destination, prevues) {
				return src
			}
			src = append(src, prevues.items[index])
			index++
		}

		prevues = nextFn(prevues)
		index = 0
	}

	return src
}

// for next := t.root; next != nil; {
// 	index, found := next.search(key)

// 	if found {
// 		return next.items[index-1].value, nil
// 	}

// 	next = next.Children[index]
// }

func CommandRange[K KeyType, V any](command string) (func(int, K, *Node[K, V]) bool, func(*Node[K, V]) *Node[K, V]) {
	switch command { // is range "<=" in relation to the source
	case "<=":
		return func(index int, destination K, node *Node[K, V]) bool {
				return node.items[index].key > destination
			}, func(n *Node[K, V]) *Node[K, V] {
				return n.nextNodeL
			}
	case ">=":
		return func(index int, destination K, node *Node[K, V]) bool {
				return node.items[index].key < destination
			}, func(n *Node[K, V]) *Node[K, V] {
				return n.nextNodeR
			}
	case "<":
		return func(index int, destination K, node *Node[K, V]) bool {
				return node.items[index].key <= destination
			}, func(n *Node[K, V]) *Node[K, V] {
				return n.nextNodeL
			}
	default:
		return func(index int, destination K, node *Node[K, V]) bool {
				return node.items[index].key >= destination
			}, func(n *Node[K, V]) *Node[K, V] {
				return n.nextNodeR
			}
	}
}
