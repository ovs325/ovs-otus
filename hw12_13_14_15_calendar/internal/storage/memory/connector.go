package memory

import (
	"context"

	ap "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/app"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"
)

type MemRepo struct {
	Repo   map[int64]st.EventModel
	LastId int64
	log    lg.Logger
}

func NewMenRepo(log lg.Logger) (ap.Storage, error) {
	repo := &MemRepo{log: log, LastId: int64(0)}
	repo.Connect(context.Background())
	return repo, nil
}

func (r *MemRepo) Connect(ctx context.Context) error {
	r.Repo = map[int64]st.EventModel{}
	return nil
}

func (r *MemRepo) Close() error {
	return nil
}
