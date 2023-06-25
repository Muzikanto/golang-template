package db

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	error2 "go-backend-template/src/utils/error"
	"go-backend-template/src/utils/logger"
	"go.uber.org/dig"
	"log"
)

// DB config

type Config interface {
	ConnString() string
}

type TxManager interface {
	RunTx(ctx context.Context, do func(ctx context.Context) error) error
}

// DB client

type DbClient struct {
	service *DbService
	pool    *pgxpool.Pool
	url     string
	ctx     context.Context
	logger  *logger.Logger
}

func CreateDB(ctx context.Context, config Config) *DbClient {
	return &DbClient{
		ctx: ctx,
		url: config.ConnString(),
	}
}

func (r *DbClient) Provide(container *dig.Scope) *DbClient {
	// provide client
	if err := container.Provide(func() *DbClient {
		return r
	}); err != nil {
		log.Fatal(err)
	}

	// provide service
	if err := container.Provide(func() *DbService {
		return r.service
	}); err != nil {
		log.Fatal(err)
	}

	// inject
	err := container.Invoke(func(logger *logger.Logger) {
		r.logger = logger
	})
	if err != nil {
		log.Fatal(err)
	}

	//
	r.logger.Log("DB initialized")

	return r
}

func (r *DbClient) Connect() *DbClient {
	r.Close()

	config, err := pgxpool.ParseConfig(r.url)

	if err != nil {
		log.Fatal(error2.Wrap(err, error2.DatabaseError, "cannot connect to database"))
	}

	pool, err := pgxpool.ConnectConfig(r.ctx, config)
	if err != nil {
		log.Fatal(error2.Wrap(err, error2.DatabaseError, "cannot connect to database"))
	}

	r.pool = pool

	//
	defer r.Close()

	r.service = CreateDbService(r)

	return r
}

func (r *DbClient) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

// DB service

type ConnManager interface {
	Conn(ctx context.Context) Connection
}

type Connection interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func CreateDbService(DbClient *DbClient) *DbService {
	return &DbService{
		DbClient: DbClient,
	}
}

type DbService struct {
	DbClient *DbClient
}

func (s *DbService) RunTx(ctx context.Context, do func(ctx context.Context) error) error {
	_, ok := hasTx(ctx)
	if ok {
		return do(ctx)
	}

	return runTx(ctx, s.DbClient, do)
}

func (s *DbService) Conn(ctx context.Context) Connection {
	tx, ok := hasTx(ctx)
	if ok {
		return tx.conn
	}

	return s.DbClient.pool
}

type txKey int

const (
	key txKey = iota
)

type transaction struct {
	conn pgx.Tx
}

func (t *transaction) commit(ctx context.Context) error {
	err := t.conn.Commit(ctx)
	if err != nil {
		return error2.Wrap(err, error2.DatabaseError, "cannot commit transaction")
	}

	return nil
}

func (t *transaction) rollback(ctx context.Context) error {
	err := t.conn.Rollback(ctx)
	if err != nil {
		return error2.Wrap(err, error2.DatabaseError, "cannot rollback transaction")
	}

	return nil
}

func withTx(ctx context.Context, tx transaction) context.Context {
	return context.WithValue(ctx, key, tx)
}

func hasTx(ctx context.Context) (transaction, bool) {
	tx, ok := ctx.Value(key).(transaction)
	if ok {
		return tx, true
	}

	return transaction{}, false
}

func runTx(ctx context.Context, DbClient *DbClient, do func(ctx context.Context) error) error {
	conn, err := DbClient.pool.Begin(ctx)
	if err != nil {
		return error2.Wrap(err, error2.DatabaseError, "cannot open transaction")
	}

	tx := transaction{conn: conn}
	txCtx := withTx(ctx, tx)

	err = do(txCtx)
	if err != nil {
		if err := tx.rollback(txCtx); err != nil {
			return err
		}
		return err
	}
	if err := tx.commit(txCtx); err != nil {
		return err
	}

	return nil
}

//var QueryBuilder = goqu.Dialect("postgres")

type Ex = goqu.Ex
type Record = goqu.Record
