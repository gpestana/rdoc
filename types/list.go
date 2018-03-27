package types

import (
	"errors"
	"fmt"
)

// A type list maintains a 2 way linked list data structure to store the data.
// It exposes an `AddElement` and `RemoveElement` which accepts a type CRDT and
// the index in which to add the new list node. The indexing starts from 0,
// which means becoming the head of the list. Index 1 is the first node after
// the head of the list, and so forth.
type List struct {
	Node
	linkedList *LinkedList
}

func NewList() *List {
	ll := LinkedList{}
	return &List{linkedList: &ll}
}

//TODO
func (l List) Delete()              {}
func (l List) AddOpPresence(string) {}
func (l List) RmOpPresence(string)  {}

func (l *List) Length() int {
	len := 0
	h := l.linkedList.Head
	if h == nil {
		return 0
	}
	len++
	for h.next != nil {
		len++
		h = h.next
	}
	return len
}

// Adds element to list (in a specific index)
// TODO: (think) is it possible to add element to index which does not exit yet?
func (l *List) AddElement(index int, el CRDT) error {
	if l.Length() < index {
		return errors.New(
			fmt.Sprintf("Index cannot be larger than the length of the list, %d", l.Length()))
	}

	if index < 0 {
		return errors.New("Index cannot be smaller than 0")
	}

	if index == 0 {
		n := LinkedListNode{
			value: el,
			next:  l.linkedList.Head,
		}
		l.linkedList.Head = &n
		currHead := l.linkedList.Head
		if currHead != nil {
			currHead.previous = &n
		}
		return nil
	}

	nPtr := l.linkedList.Head
	i := 0
	for i < index-1 {
		i++
		nPtr = nPtr.next
	}

	n := LinkedListNode{
		value:    el,
		next:     nPtr.next,
		previous: nPtr,
	}
	nPtr.next = &n
	return nil
}

// Removes element from list (in a specific index)
func (l *List) DeleteElement(index int) error {
	if l.Length() < index {
		return errors.New(
			fmt.Sprintf("Index cannot be larger than the length of the list, %d", l.Length()))
	}

	if index == 0 {
		n := l.linkedList.Head.next
		if n != nil {
			l.linkedList.Head = n
			n.previous = nil
		}
		return nil
	}

	n := l.linkedList.Head
	for 0 < index {
		n = n.next
		index--
	}

	n.previous.next = n.next
	if n.next != nil {
		n.next.previous = n.previous
	}

	return nil
}

//TODO
func (l List) String() string {
	return ""
}

type LinkedList struct {
	Head *LinkedListNode
}

func newLinkedList() *LinkedList {
	return &LinkedList{Head: nil}
}

type LinkedListNode struct {
	next     *LinkedListNode
	previous *LinkedListNode
	value    CRDT
}
