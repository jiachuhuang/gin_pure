package utils

import (
	"testing"
)

func TestQueue_Pop(t *testing.T) {
	q := NewQueue()
	q.InsertHeader(&QNode{Value:1})
	q.InsertHeader(&QNode{Value:2})
	q.InsertHeader(&QNode{Value:3})
	q.InsertHeader(&QNode{Value:4})
	q.InsertHeader(&QNode{Value:5})

	var n *QNode
	for i := 0; i < 5; i++ {
		n = q.Pop()
		if n == nil {
			t.Errorf("pop %d error", i)
		} else {
			t.Logf("pop %d ok: %d", i, n.Value.(int))
		}
	}

	if !q.Empty() {
		t.Errorf("pop error")
	}
}

func TestQueue_Shift(t *testing.T) {
	q := NewQueue()
	q.InsertHeader(&QNode{Value:1})
	q.InsertTail(&QNode{Value:11})
	q.InsertHeader(&QNode{Value:2})
	q.InsertTail(&QNode{Value:22})
	q.InsertHeader(&QNode{Value:3})
	q.InsertTail(&QNode{Value:33})
	q.InsertHeader(&QNode{Value:4})
	q.InsertTail(&QNode{Value:44})
	q.InsertHeader(&QNode{Value:5})
	q.InsertTail(&QNode{Value:55})

	var n *QNode
	for i := 0; i < 10; i++ {
		n = q.Shift()
		if n == nil {
			t.Errorf("shift %d error", i)
		} else {
			t.Logf("shift %d ok: %d", i, n.Value.(int))
		}
	}

	if !q.Empty() {
		t.Errorf("shift error")
	}
}

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	q.InsertHeader(&QNode{Value:1})
	q.InsertTail(&QNode{Value:11})
	q.InsertHeader(&QNode{Value:2})
	q.InsertTail(&QNode{Value:22})
	q.InsertHeader(&QNode{Value:3})
	q.InsertTail(&QNode{Value:33})
	q.InsertHeader(&QNode{Value:4})
	q.InsertTail(&QNode{Value:44})
	q.InsertHeader(&QNode{Value:5})
	q.InsertTail(&QNode{Value:55})

	var n *QNode
	for !q.Empty() {
		n = q.Shift()
		t.Logf("%d", n.Value.(int))
	}
}
