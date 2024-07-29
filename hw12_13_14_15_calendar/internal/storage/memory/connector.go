package memory

import (
	"context"
	"sync"

	bl "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/business_logic"
	tp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/types"
)

type MemRepo struct {
	mu     sync.RWMutex // Добавляем мьютекс
	Repo   map[int64]tp.EventModel
	LastID int64
}

func NewMemRepo() (bl.AbstractStorage, error) {
	repo := &MemRepo{LastID: int64(0)}
	repo.Connect(context.Background())
	return repo, nil
}

func (r *MemRepo) Connect(_ context.Context) error {
	r.mu.Lock()         // Блокируем мьютекс
	defer r.mu.Unlock() // Разблокируем мьютекс при выходе из функции
	r.Repo = map[int64]tp.EventModel{}
	return nil
}

func (r *MemRepo) Close() error {
	return nil
}
