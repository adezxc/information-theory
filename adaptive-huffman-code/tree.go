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

func PreGenerateTree(order int, parent *Tree) *Tree {
	if order <= 0 {
		return nil
	}

	// Create the current node
	current := NewTree(uint64(order), 1, order, nil, nil, parent, false)

	// Recursively create the left and right children
	current.Left = PreGenerateTree(order-2, current)
	current.Right = PreGenerateTree(order-1, current)

	return current
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
	t.Right = NewTree(character, 1, t.Order-1, nil, nil, t, false)
	t.Left = NewTree(0, 0, t.Order-2, nil, nil, t, true)
	t.Nyt = false
	newNode := t.Right
	nyt = t.Left

	return newNode, nyt
}

func (t *Tree) Update(mapOfTrees map[uint64]*Tree) {
	for t.Parent != nil {
		order := t.Order
		var swappableTree *Tree
		for _, tree := range mapOfTrees {
			if tree.Order > order && t.Weight+1 == tree.Weight {
				order = tree.Order
				swappableTree = tree
			}
		}
		if swappableTree != nil && swappableTree != t.Parent {
			if swappableTree.Parent.Left == swappableTree {
				if t.Parent.Left == t {
					t.Parent.Left, swappableTree.Parent.Left = swappableTree.Parent.Left, t.Parent.Left
				} else {
					t.Parent.Right, swappableTree.Parent.Left = swappableTree.Parent.Left, t.Parent.Right
				}
			} else {
				if t.Parent.Left == t {
					t.Parent.Left, swappableTree.Parent.Right = swappableTree.Parent.Right, t.Parent.Left
				} else {
					t.Parent.Right, swappableTree.Parent.Right = swappableTree.Parent.Right, t.Parent.Right
				}
			}
			t.Parent, swappableTree.Parent = swappableTree.Parent, t.Parent
			swappableTree.Order, t.Order = t.Order, swappableTree.Order
			t.GetRoot().UpdateWeights()
			t.Weight++
		} else {
			t.Weight++
		}
		t = t.Parent
	}
	t.Weight++

	return
}

func (t *Tree) UpdateWeights() {
	CopyTree := t
	if CopyTree.Left != nil && CopyTree.Right != nil {
		CopyTree.Left.UpdateWeights()
		CopyTree.Right.UpdateWeights()
		CopyTree.Weight = CopyTree.Left.Weight + CopyTree.Right.Weight
	}
}

func (t *Tree) GetRoot() *Tree {
	for t.Parent != nil {
		t = t.Parent
	}
	return t
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
