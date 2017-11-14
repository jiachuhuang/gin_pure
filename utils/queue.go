package utils

import "sync"

type Queue struct {
	root *QNode
	tail *QNode
	pool *sync.Pool
}

type QNode struct {
	prev *QNode
	next *QNode
	Value interface{}
}

func NewQueue() *Queue {
	return &Queue{nil,nil, nil}
}

func (q *Queue) NewQNode(value interface{}) *QNode {
	if q.pool == nil {
		q.pool = &sync.Pool{
			New: func() interface{} {
				return &QNode{nil,nil,nil}
			},
		}
	}
	n := q.pool.Get().(*QNode)
	n.Reset(value)
	return n
}

func (q *Queue) Recycle(n *QNode) {
	if n != nil {
		q.pool.Put(n)
	}
}

func (q *Queue) Empty() bool{
	return q.root == q.tail && q.root == nil
}

func (q *Queue) InsertHeader(n *QNode) {
	if q.Empty() {
		q.root = n
		q.tail = n
	} else {
		n.next = q.root
		q.root.prev = n
		q.root = n
	}
}

func (q *Queue) InsertTail(n *QNode) {
	if q.Empty() {
		q.root = n
		q.tail = n
	} else {
		n.prev = q.tail
		q.tail.next = n
		q.tail = n
	}
}

func (q *Queue) RemoveNode(n *QNode) {
	if n.prev != nil {
		n.prev.next = n.next
	} else if n == q.root && n != nil {
		q.root = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else if n == q.tail && n != nil {
		q.tail = n.prev
	}

	n.prev = nil
	n.next = nil
}

func (q *Queue) Pop() *QNode {
	if q.tail == nil {
		return nil
	}

	n := q.tail
	if q.tail == q.root {
		q.root = nil
		q.tail = nil
		return n
	}

	n.prev.next = nil
	q.tail = n.prev
	n.prev = nil
	return n
}

func (q *Queue) Shift() *QNode {
	if q.root == nil {
		return nil
	}

	n := q.root
	if q.tail == q.root {
		q.root = nil
		q.tail = nil
		return n
	}

	n.next.prev = nil
	q.root = n.next
	n.next = nil
	return n
}

func (q *Queue) GetTailNode() *QNode {
	return q.tail
}

func (q *Queue) GetHeaderNode() *QNode {
	return q.root
}

func (q *Queue) Clear() {
	q.root = nil
	q.tail = nil
}

func (n *QNode) Reset (value interface{}) {
	n.prev = nil
	n.next = nil
	n.Value = value
}

