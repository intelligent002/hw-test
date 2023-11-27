package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Print()
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Length int
	First  *ListItem
	Last   *ListItem
}

func (l *list) Print() {
	for item := l.Front(); item != nil; item = item.Next {
		fmt.Println(item.Value)
	}
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v}
	if l.First != nil {
		item.Next = l.First
		l.First.Prev = &item
	}
	l.First = &item

	if l.Last == nil {
		l.Last = &item
	}
	l.Length++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v}
	if l.Last != nil {
		item.Prev = l.Last
		l.Last.Next = &item
	}
	l.Last = &item

	if l.First == nil {
		l.First = &item
	}
	l.Length++
	return &item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.First == i {
		l.First = i.Next
	}
	if l.Last == i {
		l.Last = i.Prev
	}
	l.Length--
	i = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if l.First == i {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.Last == i && i.Prev != nil {
		l.Last = i.Prev
	}
	i.Prev = nil
	i.Next = l.First
	l.First = i
}

func (l *list) MoveToBack(i *ListItem) {
	if l.Last == i {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.First == i && i.Next != nil {
		l.First = i.Next
	}
	i.Next = nil
	i.Prev = l.Last
	l.Last = i
}

func NewList() List {
	return new(list)
}
