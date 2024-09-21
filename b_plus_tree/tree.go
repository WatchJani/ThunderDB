package BPTree

import (
	"bytes"
	"errors"
	"fmt"
)

type Tree[V any] struct {
	root   *Node[V]
	degree int
}

type item[V any] struct {
	key   [][]byte
	value V
}

type Node[V any] struct {
	items     []item[V]
	Children  []*Node[V]
	nextNodeL *Node[V]
	nextNodeR *Node[V]
	pointer   int
}

type Stack[V any] struct {
	store []positionStr[V]
}

type positionStr[V any] struct {
	node     *Node[V]
	position int
}

func newItem[V any](key [][]byte, value V) item[V] {
	return item[V]{
		key:   key,
		value: value,
	}
}

func newNode[V any](degree int) Node[V] {
	return Node[V]{
		items:    make([]item[V], degree+1),
		Children: make([]*Node[V], degree+2),
	}
}

func (n *Node[V]) delete() {
	n.items = nil
	n.Children = nil
	n.pointer = 0
	n.nextNodeL = nil
	n.nextNodeR = nil
}

func New[V any](degree int) *Tree[V] {
	if degree < 3 {
		degree = 3
	}

	return &Tree[V]{
		degree: degree,
		root: &Node[V]{
			items:    make([]item[V], degree+1),
			Children: make([]*Node[V], degree+2),
		},
	}
}

func (n *Node[V]) search(target [][]byte) (int, bool) {
	low, high := 0, n.pointer-1

	for low <= high {
		mid := (low + high) / 2

		var operation int = 0
		for i := 0; i < len(target); i++ {
			num := bytes.Compare(n.items[mid].key[i], target[i])

			if num != 0 {
				operation = num
				break
			}
		}

		if operation == 0 {
			return mid + 1, true
		} else if operation == -1 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low, false
}

func newStack[V any]() Stack[V] {
	return Stack[V]{
		store: make([]positionStr[V], 0, 4),
	}
}

func (s *Stack[V]) Push(node *Node[V], position int) {
	s.store = append(s.store, positionStr[V]{
		node:     node,
		position: position,
	})
}

func (s *Stack[V]) Pop() (positionStr[V], error) {
	if len(s.store) == 0 {
		return positionStr[V]{}, errors.New("stack is empty")
	}

	pop := s.store[len(s.store)-1]
	s.store = s.store[:len(s.store)-1]
	return pop, nil
}

func (t *Tree[V]) Find(key [][]byte, operation string) (*Node[V], int, error) {
	var (
		prevues  *Node[V]
		res      int = -1
		position int
	)

	for current := t.root; current != nil; {
		position, _ = current.search(key)

		prevues = current
		current = current.Children[position]
	}

	if operation == ">" {
		res = position
		if prevues.pointer == position {
			res = -1
		}
	} else if operation == "<" {
		if position >= 1 {
			res = position - 1
		} else {
			res = -1
		}
	}

	return prevues, res, fmt.Errorf("key %v not found", key)
}

func (t *Node[V]) GetItem(index int) item[V] {
	return t.items[index]
}

func (t *Node[V]) GetValue(index int) V {
	return t.items[index].value
}

func (t *Node[V]) NextLeft() *Node[V] {
	return t.nextNodeL
}

func (t *Node[V]) NextRight() *Node[V] {
	return t.nextNodeR
}

func (t *Tree[V]) Insert(key [][]byte, value V) {
	stack, item := newStack[V](), newItem(key, value)
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

			newNode := newNode[V](t.degree)
			newNode.pointer += migrate(newNode.items, parent.items[:middle], 0) //migrate half element to left child node
			migrate(newNode.Children, parent.Children[:middle+1], 0)

			parent.pointer -= deleteElement(parent.items, 0, parent.pointer-middle+1-t.degree&1)
			migrate(parent.Children, parent.Children[middle+1:], 0)
			nodeChildren = &newNode

			current = stack
		}

		rootNode := newNode[V](t.degree)
		rootNode.pointer += insert(rootNode.items, middleKey, 0)
		rootNode.Children[0] = nodeChildren
		rootNode.Children[1] = current.node
		t.root = &rootNode
	}
}

func childrenIndex(key, value [][]byte, index int) int {
	for i := 0; i < len(key); i++ {
		num := bytes.Compare(value[i], key[i])
		if num == -1 {
			return index + 1
		} else if num == 1 {
			return index
		}
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

func insertLeaf[V any](current *Node[V], position, degree int, item item[V]) (item[V], *Node[V]) {
	current.pointer += insert(current.items, item, position)

	if current.pointer < degree {
		return item, nil
	}

	newNode := newNode[V](degree)
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

func findLeaf[V any](root *Node[V], stack *Stack[V], key [][]byte) (int, bool) {
	position, found := 0, false

	for current := root; current != nil; {
		position, found = current.search(key)
		stack.Push(current, position)

		current = current.Children[position]
	}

	return position, found
}

func (t *Tree[V]) Delete(key [][]byte) error {
	stack := newStack[V]()
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

func siblingExist[V any](parent positionStr[V], index int) (*Node[V], bool) {
	index = parent.position + index

	if index < 0 || index > parent.node.pointer {
		return nil, false
	}

	sibling := parent.node.Children[index]
	return sibling, true
}

func sibling[V any](parent positionStr[V], degree int) (*Node[V], bool, bool) {
	var (
		potential *Node[V]
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

func merge[V any](current, sibling *Node[V], parent positionStr[V], leafInternal, side bool) {
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

func transfer[V any](parent, current positionStr[V], sibling *Node[V], leafInternal, side bool) {
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

func (t *Tree[V]) TestFunc() int {
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

func (t *Tree[V]) GetRoot() *Node[V] {
	return t.root
}

func (t *Tree[V]) BetweenKey(key [][]byte) (V, error) {
	next, index := t.root, 0
	var closestValue V

	if t.root.pointer == 0 {
		return closestValue, fmt.Errorf("key %v not found and no suitable larger or smaller key exists", key)
	}

	var privies *Node[V]
	for next != nil {
		index, _ = next.search(key)
		privies = next

		next = next.Children[index]
	}

	if index > 0 {
		index--
	}

	return privies.items[index].value, nil
}
