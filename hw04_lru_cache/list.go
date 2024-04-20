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
	List []*ListItem
}

func NewList() List {
	return new(list)
}

// длина списка.
func (l *list) Len() int {
	return len(l.List)
}

// первый элемент списка.
func (l *list) Front() *ListItem {
	if l.Len() > 0 {
		return l.List[0]
	}
	return nil
}

// последний элемент списка.
func (l *list) Back() *ListItem {
	if l.Len() > 0 {
		return l.List[l.Len()-1]
	}
	return nil
}

// добавить значение в начало.
func (l *list) PushFront(item any) *ListItem {
	if item == nil {
		return nil
	}
	newElement := ListItem{Value: item}
	if l.Len() == 0 {
		l.List = []*ListItem{&newElement}
		return &newElement
	}
	newElement.Next = l.Front()
	l.Front().Prev = &newElement
	newListAny := make([]*ListItem, l.Len()+1)
	copy(newListAny[1:], l.List)
	newListAny[0] = &newElement
	l.List = newListAny
	return &newElement
}

// добавить значение в конец.
func (l *list) PushBack(item any) *ListItem {
	if item == nil {
		return nil
	}
	newElement := ListItem{
		Value: item,
		Next:  nil,
		Prev:  l.Back(),
	}
	l.Back().Next = &newElement
	l.List = append(l.List, &newElement)
	return &newElement
}

// удалить элемент.
func (l *list) Remove(itemRm *ListItem) {
	switch itemRm {
	case nil:
		return
	case l.Front():
		l.List = l.List[1:l.Len()]
		l.Front().Prev = nil
		return
	case l.Back():
		l.List = l.List[0 : l.Len()-1]
		l.Back().Next = nil
		return
	}
	for i, item := range l.List {
		if item.Next == itemRm {
			listTail := make([]*ListItem, l.Len()-i-2)
			copy(listTail, l.List[i+2:l.Len()])
			l.List = l.List[0 : i+1]
			l.Back().Next = listTail[0]
			listTail[0].Prev = l.Back()
			l.List = append(l.List, listTail...)
			break
		}
	}
}

// переместить элемент в начало.
func (l *list) MoveToFront(itemMv *ListItem) {
	if itemMv == nil || l.Front() == itemMv {
		return
	}
	l.Remove(itemMv)
	l.PushFront(itemMv.Value)
}
