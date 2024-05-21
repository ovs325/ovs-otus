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
	firstItem *ListItem
	lastItem  *ListItem
	lenList   int
}

func NewList() List {
	return new(list)
}

// длина списка.
func (l *list) Len() (num int) {
	return l.lenList
}

// первый элемент списка.
func (l *list) Front() *ListItem {
	return l.firstItem
}

// последний элемент списка.
func (l *list) Back() *ListItem {
	return l.lastItem
}

// добавить значение в начало.
func (l *list) PushFront(item any) *ListItem {
	if item == nil {
		return nil
	}
	newElement := ListItem{Value: item}
	if l.lenList == 0 {
		l.firstItem = &newElement
		l.lastItem = &newElement
	} else {
		l.firstItem.Prev = &newElement
		newElement.Next = l.firstItem
		l.firstItem = &newElement
	}
	l.lenList++
	return l.firstItem
}

// добавить значение в конец.
func (l *list) PushBack(item any) *ListItem {
	if l.firstItem == nil {
		l.PushFront(item)
		return l.firstItem
	}
	newElement := ListItem{Value: item}
	if l.lenList == 0 {
		l.firstItem = &newElement
		l.lastItem = &newElement
	} else {
		l.lastItem.Next = &newElement
		newElement.Prev = l.lastItem
		l.lastItem = &newElement
	}
	l.lenList++
	return l.lastItem
}

// удалить элемент.
func (l *list) Remove(itemRm *ListItem) {
	switch itemRm {
	case nil:
		return
	case l.firstItem:
		l.firstItem = l.firstItem.Next
		if l.firstItem != nil {
			l.firstItem.Prev = nil
		}
	case l.lastItem:
		l.lastItem = l.lastItem.Prev
		l.lastItem.Next = nil
	default:
		itemRm.Prev.Next = itemRm.Next
		itemRm.Next.Prev = itemRm.Prev
	}
	l.lenList--
}

// переместить элемент в начало.
func (l *list) MoveToFront(itemMv *ListItem) {
	switch itemMv {
	case nil, l.firstItem:
		return
	default:
		l.Remove(itemMv)
		itemMv.Prev = nil
		l.firstItem.Prev = itemMv
		itemMv.Next = l.firstItem
		l.firstItem = itemMv
		l.lenList++
	}
}
