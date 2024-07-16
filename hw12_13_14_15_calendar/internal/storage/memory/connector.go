package memory

import (
	"context"
	"fmt"

	ap "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/app"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	st "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage"

	rx "github.com/restream/reindexer"
)

type RxRepo struct {
	DB        *rx.Reindexer
	DSN       string
	Namespace string
	log       lg.Logger
}

func NewRxRepo(ctx context.Context, conf *cf.Config, log lg.Logger) (ap.Storage, error) {
	repo := &RxRepo{Namespace: conf.RxCnf.Namespace}
	repo.GetDSN(conf)
	if err := repo.Connect(ctx); err != nil {
		return nil, err
	}
	if err := repo.CreateTables(ctx); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RxRepo) Connect(ctx context.Context) error {

	r.log.Info("ConnectDB started")

	r.DB = rx.NewReindex(r.DSN, rx.WithCreateDBIfMissing())

	if !r.ping() {
		return fmt.Errorf("failed to connect to db")
	}
	r.log.Info("The connection has been successfully established")

	return nil
}

func (r *RxRepo) Close() error {
	r.log.Info("Closing database connection")
	if r.DB != nil {
		r.DB.Close()
	}
	r.log.Info("Database connection closed successfully")
	return nil
}

func (r *RxRepo) CreateTables(ctx context.Context) (err error) {
	r.log.Info("Creating namespace", "name", r.Namespace)

	namespaces, err := r.DB.DescribeNamespaces()
	if err != nil {
		return fmt.Errorf("failed to get the list of 'namespaces': %w", err)
	}
	isExist := false
	for _, namespaceItem := range namespaces {
		if namespaceItem.Name == r.Namespace {
			if err := r.DB.TruncateNamespace(r.Namespace); err != nil {
				return fmt.Errorf("could not clear 'namespace' %s: %w", r.Namespace, err)
			}
			isExist = true
			break
		}
	}
	if !isExist {
		err := r.DB.OpenNamespace(r.Namespace, rx.DefaultNamespaceOptions(), st.EventModel{})
		if err != nil {
			return fmt.Errorf("failed to create 'namespace' %s: %w", r.Namespace, err)
		}
	}
	r.log.Info("The creating namespace has been successfully established")
	return nil
}

func (r *RxRepo) GetDSN(d *cf.Config) string {
	r.DSN = fmt.Sprintf(
		"cproto://%v:%v/%v",
		d.RxCnf.Host,
		d.RxCnf.Port,
		d.Db.Database,
	)
	return r.DSN
}

func (r *RxRepo) ping() bool {

	type PingItem struct {
		ID   int    `reindex:"id,tree,pk"`
		Name string `reindex:"name"`
	}

	ns := "ping"

	namespaces, err := r.DB.DescribeNamespaces()
	if err != nil {
		return false
	}
	isExist := false
	for _, namespaceItem := range namespaces {
		if namespaceItem.Name == ns {
			if err := r.DB.TruncateNamespace(ns); err != nil {
				return false
			}
			isExist = true
			break
		}
	}
	if !isExist {
		r.DB.OpenNamespace(ns, rx.DefaultNamespaceOptions(), PingItem{})
	}
	defer func() { _ = r.DB.CloseNamespace(ns) }()

	id := 1

	if err := r.DB.Upsert(ns, &PingItem{ID: id, Name: "Ping"}); err != nil {
		return false
	}

	if _, found := r.DB.Query(ns).Where("id", rx.EQ, id).Get(); !found {
		return false
	}
	return true
}
