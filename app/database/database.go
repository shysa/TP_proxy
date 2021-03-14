package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/shysa/TP_proxy/config"
	"github.com/sirupsen/logrus"
	"os"
)

type DB struct {
	dbPool *pgxpool.Pool
	config *config.ConfDB
}

func NewDB(config *config.ConfDB) *DB {
	return &DB{
		config: config,
	}
}

func (db *DB) Open() error {
	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=%s pool_max_conns=%s",
		db.config.Username,
		db.config.Password,
		db.config.Host,
		db.config.DbName,
		db.config.SslMode,
		db.config.MaxConn,
	))
	if err != nil {
		return err
	}

	looger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.ErrorLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	conf.ConnConfig.Logger = logrusadapter.NewLogger(looger)

	db.dbPool, err = pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() {
	db.dbPool.Close()
}

func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.dbPool.Begin(ctx)
}

func (db *DB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return db.dbPool.Exec(ctx, sql, arguments...)
}

func (db *DB) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error) {
	return db.dbPool.Query(ctx, sql, optionsAndArgs...)
}

func (db *DB) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
	return db.dbPool.QueryRow(ctx, sql, optionsAndArgs...)
}

func (db *DB) CopyFrom(ctx context.Context, table pgx.Identifier, cols []string, src pgx.CopyFromSource) (int64, error) {
	return db.dbPool.CopyFrom(ctx, table, cols, src)
}