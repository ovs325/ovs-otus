package sql

// /home/ovs/90-Involta/500 GitLab/manual version/srv/notifications/db/init-database/init.sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	ap "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/app"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
)

type PgRepo struct {
	DB  *pgxpool.Pool //nolint:gochecknoglobals
	DSN string
	log lg.Logger
}

func NewPgRepo(ctx context.Context, conf *cf.Config, log lg.Logger) (ap.Storage, error) {
	repo := &PgRepo{log: log}
	repo.GetDSN(conf)
	if err := repo.Connect(ctx); err != nil {
		return nil, err
	}
	if err := repo.CreateTables(ctx); err != nil {
		return nil, err
	}
	return repo, nil
}

func (s *PgRepo) Connect(ctx context.Context) error {
	s.log.Info("ConnectDB started")
	conn, err := pgxpool.Connect(ctx, s.DSN)
	if err != nil {
		return fmt.Errorf("the connection attempt failed: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	s.log.Info("The connection has been successfully established")

	s.DB = conn
	return nil
}

func (s *PgRepo) Close() error { //nolint:stylecheck
	s.log.Info("Closing database connection")
	if s.DB != nil {
		s.DB.Close()
	}
	s.log.Info("Database connection closed successfully")
	return nil
}

func (s *PgRepo) CreateTables(ctx context.Context) (err error) {
	s.log.Info("Creating tables")
	for _, item := range s.GetSQLs1() {
		if _, err = s.DB.Exec(ctx, item[0]); err != nil {
			return fmt.Errorf("%s %w", item[1], err)
		}
	}
	// Check if table was created
	sql := `
SELECT COUNT(*) 
FROM information_schema.tables 
WHERE table_name = 'events'`
	var count int
	if err = s.DB.QueryRow(ctx, sql).Scan(&count); err != nil {
		return fmt.Errorf("failed to check table existence. %w", err)
	}
	if count == 0 {
		return fmt.Errorf("table 'events' was not created")
	}
	s.log.Info("Table 'events' created successfully")

	for _, item := range s.GetSQLs2() {
		if _, err = s.DB.Exec(ctx, item[0]); err != nil {
			return fmt.Errorf("%s %w", item[1], err)
		}
	}

	s.log.Info("Tables created successfully")
	return nil
}

func (s *PgRepo) GetDSN(d *cf.Config) string {
	s.DSN = fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		d.PgCnf.User,
		d.PgCnf.Password,
		d.PgCnf.Host,
		d.PgCnf.Port,
		d.DB.Database,
	)
	return s.DSN
}
