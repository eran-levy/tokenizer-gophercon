package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/repository"
	"github.com/eran-levy/tokenizer-gophercon/repository/model"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type mysqlRepository struct {
	config    repository.Config
	db        *sql.DB
	telemetry telemetry.Telemetry
}

const metadataTableName = "tokenizer_metadata"

var sqlStmntTagKey = label.Key("db.statement")

func (m mysqlRepository) StoreMetadata(ctx context.Context, mtd model.TokenizeTextMetadata) error {
	ctx, span := m.telemetry.Tracer.Start(ctx, "store tokenizer metadata")
	defer span.End()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		span.SetStatus(codes.Error, "store metadata failed, cant begin tx")
		return errors.Wrap(err, "store metadata failed, cant begin tx") // proper error handling instead of panic in your app
	}
	defer tx.Commit()
	p := fmt.Sprintf("INSERT INTO %s VALUES( ?, ?, ?, ? )", metadataTableName)
	span.SetAttributes(sqlStmntTagKey.String(p))
	stmtIns, err := m.db.Prepare(p) // ? = placeholder
	if err != nil {
		span.SetStatus(codes.Error, "store metadata failed")
		return errors.Wrap(err, "store metadata failed") // proper error handling instead of panic in your app
	}
	_, err = stmtIns.ExecContext(ctx, mtd.RequestId, mtd.GlobalTxId, mtd.CreatedDate, mtd.Language)
	if err != nil {
		span.SetStatus(codes.Error, "store metadata failed")
		return errors.Wrap(err, "store metadata failed") // proper error handling instead of panic in your app
	}
	err = stmtIns.Close()
	if err != nil {
		span.SetStatus(codes.Error, "store metadata failed")
		return errors.Wrap(err, "store metadata failed")
	}
	return nil
}

func New(config repository.Config, telemetry telemetry.Telemetry) (repository.Persistence, error) {
	//dsn - https://github.com/go-sql-driver/mysql#dsn-data-source-name
	//example user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
	//https://github.com/go-sql-driver/mysql#usage
	const dbDriver = "mysql"
	db, err := sql.Open(dbDriver, config.Dsn)
	if err != nil {
		return mysqlRepository{}, errors.Wrap(err, "could not connect to database")
	}
	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	err = db.Ping()
	if err != nil {
		return mysqlRepository{}, err
	}
	return mysqlRepository{db: db, telemetry: telemetry}, nil
}

func (m mysqlRepository) Close() error {
	err := m.db.Close()
	if err != nil {
		return err
	}
	logger.Log.Info("gracefully closed repo client")
	return nil
}

func (m mysqlRepository) IsServiceHealthy(ctx context.Context) (bool, error) {
	err := m.db.Ping()
	if err != nil {
		return false, errors.Wrap(err, "db ping repond with an error")
	}
	return true, nil
}
