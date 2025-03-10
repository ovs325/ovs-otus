package listsystem

import (
	"sync"

	tp "bruteforce/internal/types"
)

type SpecList struct {
	lock         sync.Mutex // Для потокобезопастности
	nameSpecList string
	idsMap       map[string]tp.IsBlack // для подсетей ключи типа: "192.1.1.0/25", "10.0.0.0/8"
	idSeeker     func(subnetMap map[string]tp.IsBlack, ipStr string) (isFound bool, isBlack tp.IsBlack, err error)
	keyVerifier  func(subnetStr string) (ok bool)
}

func NewSpecList(
	name string,
	lenMap int,
	seeker func(subnetMap map[string]tp.IsBlack, ipStr string) (isFound bool, isBlack tp.IsBlack, err error),
	verifier func(subnetStr string) (ok bool),
) *SpecList {
	sl := SpecList{
		nameSpecList: name,
		idSeeker:     seeker,
		keyVerifier:  verifier,
	}
	if lenMap > 0 {
		sl.idsMap = make(map[string]tp.IsBlack, lenMap)
	} else {
		sl.idsMap = map[string]tp.IsBlack{}
	}
	return &sl
}

func (s *SpecList) ToFind(id string) (isFound bool, isBlack tp.IsBlack, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.idSeeker(s.idsMap, id)
}

func (s *SpecList) AddId(id string, isBlack tp.IsBlack) (ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.keyVerifier(id) {
		return false
	}
	s.idsMap[id] = isBlack
	val, ok := s.idsMap[id]
	return ok && val == isBlack
}

func (s *SpecList) DelId(id string) (ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.idsMap, id)
	_, ok = s.idsMap[id]
	return !ok
}

func (s *SpecList) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.idsMap = map[string]tp.IsBlack{}
}

func (s *SpecList) GetName() string {
	return s.nameSpecList
}

func (s *SpecList) GetIdsMap() map[string]tp.IsBlack {
	return s.idsMap
}
