package pqueue

import (
	"fmt"
	"strings"
)

type PQueue[Node comparable] struct {
	Head *PQueueNode[Node]
}

type PQueueNode[Node comparable] struct {
	Node     Node
	Priority int
	Next     *PQueueNode[Node]
}

func (pq *PQueue[Node]) String() string {
	var ret strings.Builder
	cursor := pq.Head
	for cursor != nil {
		ret.WriteString(fmt.Sprintf("%v=%d, ", cursor.Node, cursor.Priority))
		cursor = cursor.Next
	}
	return ret.String()
}
func (pq *PQueue[Node]) Print() {
	cursor := pq.Head
	for cursor != nil {
		fmt.Println(cursor.Node, " ", cursor.Priority)
		cursor = cursor.Next
	}
}

func (pq *PQueue[Node]) Pop() Node {
	if pq.Head == nil {
		panic("pop on empty pqueue")
	}
	ret := pq.Head.Node
	pq.Head = pq.Head.Next
	return ret
}

func (pq *PQueue[Node]) AddWithPriority(node Node, prio int) {
	newnode := &PQueueNode[Node]{
		Node:     node,
		Priority: prio,
	}

	if pq.Head == nil {
		pq.Head = newnode
		return
	}

	if pq.Head.Priority > prio {
		newnode.Next = pq.Head
		pq.Head = newnode
		return
	}

	cursor := pq.Head
	for {
		if cursor.Next == nil || cursor.Next.Priority > prio {
			break
		}
		cursor = cursor.Next
	}

	newnode.Next, cursor.Next = cursor.Next, newnode
	return
}

func (pq *PQueue[Node]) Has(c Node) bool {
	n := pq.Head
	for n != nil && n.Node != c {
		n = n.Next
	}
	return n != nil && n.Node == c
}
