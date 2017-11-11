package utils

type Queue struct {
	root *QNode
	tail *QNode
}

type QNode struct {
	prev *QNode
	next *QNode
	value interface{}
}

func NewQueue() *Queue {
	return &Queue{nil,nil}
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
	}

	if n.next != nil {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = nil
}

