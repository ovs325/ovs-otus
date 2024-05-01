package hw04lrucache

import "fmt"

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
	newListAny := make([]*ListItem, 0)
	newListAny = append(newListAny, &newElement)
	newListAny = append(newListAny, l.List...)
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

// // удалить элемент.
// func (l *list) Remove(itemRm *ListItem) {
// 	switch itemRm {
// 	case nil:
// 		return
// 	case l.Back():
// 		l.List = l.List[0 : l.Len()-1]
// 		l.Back().Next = nil
// 		return
// 	}
// 	newListAny := make([]*ListItem, 0)
// 	for i, item := range l.List {
// 		if item == itemRm {
// 			if item.Next == nil {
// 				l.List = newListAny
// 				break
// 			}
// 			l.List[i+1].Prev = itemRm.Prev
// 			newListAny = append(newListAny, l.List[i+1:l.Len()]...)
// 			l.List = newListAny
// 			break
// 		} else {
// 			if item.Next == itemRm {
// 				item.Next = itemRm.Next
// 			}
// 			newListAny = append(newListAny, item)
// 		}
// 	}
// }

// переместить элемент в начало.
func (l *list) MoveToFront(itemMv *ListItem) {
	if itemMv == nil || l.Front() == itemMv {
		return
	}
	l.Remove(itemMv)
	l.PushFront(itemMv.Value)
}

// удалить элемент.
func (l *list) Remove(itemRm *ListItem) {
	if itemRm == nil {
		return
	}

	num := 0
	ok := GetNum(&num, itemRm, l.Front())

	if !ok {
		fmt.Println("Не найдено!")
	}

	newListAny := make([]*ListItem, 0)
	newListAny = append(newListAny, l.List[0:num]...)
	newListAny[num-1].Next = itemRm.Next
	if itemRm.Next != nil {
		newListAny = append(newListAny, l.List[num+1:l.Len()]...)
		newListAny[num].Prev = itemRm.Prev
	}
	l.List = newListAny
}

// Функция для поиска места элемента по цепочке ссылок.
func GetNum(numPoint *int, delItem, item *ListItem) bool {
	if delItem == item {
		return true
	}
	*numPoint++
	if item.Next == nil {
		return false
	}
	return GetNum(numPoint, delItem, item.Next)
}
