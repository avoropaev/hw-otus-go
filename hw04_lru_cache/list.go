package hw04lrucache

type List interface {
	Len() int
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
	Length    int
	FrontItem *ListItem
	BackItem  *ListItem
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.FrontItem
}

func (l *list) Back() *ListItem {
	return l.BackItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.FrontItem != nil {
		l.FrontItem.Prev = newItem
		newItem.Next = l.FrontItem
	}

	l.FrontItem = newItem

	if l.BackItem == nil {
		l.BackItem = newItem
	}

	l.Length++

	return l.FrontItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.BackItem != nil {
		l.BackItem.Next = newItem
		newItem.Prev = l.BackItem
	}

	l.BackItem = newItem

	if l.FrontItem == nil {
		l.FrontItem = newItem
	}

	l.Length++

	return l.BackItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if l.FrontItem == i {
		l.FrontItem = i.Next
	}

	if l.BackItem == i {
		l.BackItem = i.Prev
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	l.Remove(i)

	if l.FrontItem != nil {
		l.FrontItem.Prev = i
		i.Next = l.FrontItem
	}

	l.FrontItem = i

	if l.BackItem == nil {
		l.BackItem = i
	}

	l.Length++
}

func NewList() List {
	return new(list)
}
