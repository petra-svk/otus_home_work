package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	total int
}

func (l *list) Len() int {
	return l.total
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{Value: v}
	if l.first == nil {
		l.first = &newItem
		l.last = &newItem
	} else {
		newItem.Next = l.first
		l.first.Prev = &newItem
		l.first = &newItem
	}
	l.total++
	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.last == nil {
		return l.PushFront(v)
	}

	newItem := ListItem{Value: v}
	l.last.Next = &newItem
	newItem.Prev = l.last
	l.last = &newItem
	l.total++
	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.total--
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case l.Front():
		return
	case l.Back():
		i.Prev.Next = nil
		l.last = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func (l *list) Clear() {
	l.first = nil
	l.last = nil
	l.total = 0
}

func NewList() List {
	return new(list)
}
