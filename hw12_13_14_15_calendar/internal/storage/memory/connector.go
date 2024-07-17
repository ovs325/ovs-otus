package memory

import (
	"context"
	"sync"

	ap "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/app"
	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"
)

type MemRepo struct {
	mu     sync.RWMutex // Добавляем мьютекс
	Repo   map[int64]st.EventModel
	LastID int64
}

func NewMemRepo() (ap.Storage, error) {
	repo := &MemRepo{LastID: int64(0)}
	repo.Connect(context.Background())
	return repo, nil
}

func (r *MemRepo) Connect(_ context.Context) error {
	r.mu.Lock()         // Блокируем мьютекс
	defer r.mu.Unlock() // Разблокируем мьютекс при выходе из функции
	r.Repo = map[int64]st.EventModel{}
	return nil
}

func (r *MemRepo) Close() error {
	return nil
}
