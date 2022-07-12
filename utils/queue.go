package utils

import "strings"

// Queue 优先级队列实现
type Queue interface {
	GetSize() int
	IsEmpty() bool
	Push(node Node)
	Pop() *Node
	Top() *Node
}

type Node struct {
	Point int    		// 分数
	Val   interface{} 	// 值
}

// PriorityQueue 按照分数大小排列
type PriorityQueue struct {
	Size int
	Item []Node
}

func (p *PriorityQueue) GetSize() int {
	return p.Size
}

func (p *PriorityQueue) IsEmpty() bool {
	return len(p.Item) == 0
}

func (p *PriorityQueue) Push(node Node) {

	i := p.Size
	p.Size++

	for {
		if i <= 0 {
			break
		}
		parent := (i - 1) / 2
		if p.Item[parent].Point >= node.Point {
			break
		}
		p.Item[parent], p.Item[i] = p.Item[i], p.Item[parent]
		i = parent
	}

	p.Item[i] = node
}

func (p *PriorityQueue) Pop() *Node {

	if p.Size == 0 {
		return nil
	}
	root := p.Item[0]
	i := 0
	p.Size--

	last := p.Item[p.Size]
	p.Item[p.Size] = Node{}

	for {
		left := 2*i + 1
		right := 2*i + 2

		if left >= p.Size {
			break
		}

		if right > p.Size && p.Item[left].Point < p.Item[right].Point {
			left = right
		}

		if p.Item[left].Point < last.Point {
			break
		}

		p.Item[i], p.Item[left] = p.Item[left], p.Item[i]
		i = left
	}

	p.Item[i] = last
	return &root
}

func (p *PriorityQueue) Top() *Node {
	return &p.Item[0]
}

func NewPriorityQueue() Queue {
	return &PriorityQueue{Size: 0, Item: make([]Node, 16)}
}


type FifoQueue struct {
	Item []Node
}

func (fq *FifoQueue) GetSize() int {
	return len(fq.Item)
}

func (fq *FifoQueue) IsEmpty() bool{
	return len(fq.Item) == 0
}

func (fq *FifoQueue) Push(node Node){
	fq.Item = append(fq.Item, node)
}

func (fq *FifoQueue) Pop() *Node{
	item := &fq.Item[0]
	fq.Item = fq.Item[1:]
	return item
}

func (fq *FifoQueue) Top() *Node{
	return &fq.Item[0]
}

func NewFifoQueue() Queue {
	return &FifoQueue{Item: []Node{}}
}


type LifoQueue struct {
	FifoQueue
}

func (lq *LifoQueue) Pop() *Node{
	item := &lq.Item[lq.GetSize()-1]
	lq.Item = lq.Item[:lq.GetSize()-1]
	return item
}

func NewLifoQueue() Queue {
	return new(LifoQueue)
}

const (
	Priority = "PRIORITY"
	LIFO = "LIFO"
	FIFO = "FIFO"
)
// NewQueue 根据队列类型名称返回对应的的队列实例
// 默认返回 FIFO 队列
func NewQueue(queueType string) Queue {
	switch strings.ToUpper(queueType) {
	case Priority:
		return NewPriorityQueue()
	case LIFO:
		return NewLifoQueue()
	case FIFO:
		return NewFifoQueue()
	default:
		return NewFifoQueue()
	}

}