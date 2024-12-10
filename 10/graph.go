package main

type direction int

const (
	north direction = iota
	east
	south
	west
)

type node struct {
	score        int
	height       int
	attached     []*node
	destinations []*node
}

func newNode(height int) *node {
	attached := make([]*node, 4, 4)
	return &node{
		score:        0,
		height:       height,
		attached:     attached,
		destinations: []*node{},
	}
}

func (n *node) attach(dir direction, next *node) {
	n.attached[dir] = next
	next.attached[oposite(dir)] = n
}

func (n *node) backtrack(visitor func(n *node)) {
	for dir := range 4 {
		next := n.attached[dir]
		if next == nil || next.height != n.height-1 {
			continue
		}
		visitor(next)
		if next.height == start {
			continue
		}
		next.backtrack(visitor)
	}
}

func oposite(dir direction) direction {
	switch dir {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}
	return north
}
