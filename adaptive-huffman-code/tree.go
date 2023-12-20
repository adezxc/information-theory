package main

import (
	"math"
)

type Tree struct {
	Value  uint64
	Weight int
	Order  int

	Nyt bool

	Left   *Tree
	Right  *Tree
	Parent *Tree
}

func EmptyTree() *Tree {
	return &Tree{}
}

func NewTree(value uint64, weight, order int, left, right, parent *Tree, nyt bool) *Tree {
	return &Tree{
		Value:  value,
		Weight: weight,
		Order:  order,
		Left:   left,
		Right:  right,
		Parent: parent,
		Nyt:    nyt,
	}
}

func (t *Tree) ProcessNewCharacter(character uint64) (newChar, nyt *Tree) {
	t.Right = NewTree(character, 0, t.Order-1, nil, nil, t, false)
	t.Left = NewTree(0, 0, t.Order-2, nil, nil, t, true)
	t.Nyt = false

	newNode := t.Right
	nyt = t.Left
	newNode.Update()

	return newNode, nyt
}

func (t *Tree) Update() {
	for t.Parent != nil {
		var otherNode *Tree
		if t != t.Parent.Right {
			otherNode = t.Parent.Right
		} else {
			otherNode = t.Parent.Left
		}
		if t.Weight+1 > otherNode.Weight && !otherNode.Nyt {
			t.Parent.Right, t.Parent.Left = t.Parent.Left, t.Parent.Right
			t.Parent.Right.Order, t.Parent.Left.Order = t.Parent.Left.Order, t.Parent.Right.Order
		}
		t.Weight++
		t = t.Parent
	}
	t.Weight++
}

func (t *Tree) GetTreeIndex() (result, index uint64) {
	treePointerCopy := t

	var bitSet int
	for treePointerCopy.Parent != nil {
		if treePointerCopy.Order > treePointerCopy.Parent.Left.Order {
			bitSet = 1
		} else {
			bitSet = 0
		}
		result += uint64(bitSet * int(math.Pow(2, float64(index))))
		index++
		treePointerCopy = treePointerCopy.Parent
	}

	return result, index
}
