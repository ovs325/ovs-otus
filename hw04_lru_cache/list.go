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
}

func NewList() List {
	return new(list)
}

// длина списка.
func (l *list) Len() (num int) {
	if l.firstItem == nil {
		return 0
	}
	num = 1
	GetLen(&num, l.firstItem)
	return
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
	if l.firstItem == nil {
		l.firstItem = &newElement
		return l.firstItem
	}
	newElement.Next = l.firstItem
	l.firstItem.Prev = &newElement
	if l.lastItem == nil {
		l.lastItem = l.firstItem
	}
	l.firstItem = &newElement
	if l.lastItem == nil {
		l.firstItem = l.firstItem.Next
	}
	return l.firstItem
}

// добавить значение в конец.
func (l *list) PushBack(item any) *ListItem {
	if l.firstItem == nil {
		l.PushFront(item)
		return l.firstItem
	}
	newElement := ListItem{
		Value: item,
	}
	if l.lastItem == nil {
		newElement.Prev = l.firstItem
		l.firstItem.Next = &newElement
		l.lastItem = &newElement
		return l.lastItem
	}

	newElement.Prev = l.lastItem
	l.lastItem.Next = &newElement
	l.lastItem = &newElement
	return l.lastItem
}

// удалить элемент.
func (l *list) Remove(itemRm *ListItem) {
	switch itemRm {
	case nil:
		return
	case l.firstItem:
		l.firstItem = l.firstItem.Next
		l.firstItem.Prev = nil
		if l.firstItem.Next == nil {
			l.lastItem = nil
		}
	case l.lastItem:
		l.lastItem = l.lastItem.Prev
		l.lastItem.Next = nil
		if l.lastItem.Prev == nil {
			l.lastItem = nil
		}
	default:
		itemRm.Prev.Next = itemRm.Next
		itemRm.Next.Prev = itemRm.Prev
	}
}

// переместить элемент в начало.
func (l *list) MoveToFront(itemMv *ListItem) {
	if itemMv == nil || l.Front() == itemMv {
		return
	}
	l.Remove(itemMv)
	itemMv.Next = l.firstItem
	itemMv.Prev = nil
	l.firstItem.Prev = itemMv
	l.firstItem = itemMv
	item := GetLast(l.firstItem)
	if l.firstItem == item {
		l.lastItem = nil
	} else {
		l.lastItem = item
	}
}

func GetLen(numPoint *int, item *ListItem) {
	if item.Next == nil {
		return
	}
	*numPoint++
	GetLen(numPoint, item.Next)
}

func GetLast(item *ListItem) (last *ListItem) {
	if item.Next == nil {
		return item
	}
	return GetLast(item.Next)
}
