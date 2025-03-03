package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var (
	StatementBuilder  = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	DefaultMaxRetries = 3
)

func QueryContext(ctx context.Context, builder sq.Sqlizer) (*sql.Rows, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("Error building SQL: %v", err)
		return nil, err
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	var rows *sql.Rows
	operation := func() error {
		var opErr error
		rows, opErr = DB.QueryContext(ctx, query, args...)
		return opErr
	}

	err = RetryOperation(operation, DefaultMaxRetries)
	return rows, err
}

func QueryRowContext(ctx context.Context, builder sq.Sqlizer) *sql.Row {
	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("Error building SQL: %v", err)
		return nil
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return DB.QueryRowContext(ctx, query, args...)
}

func ExecContext(ctx context.Context, builder sq.Sqlizer) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("Error building SQL: %v", err)
		return nil, err
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	var result sql.Result
	operation := func() error {
		var opErr error
		result, opErr = DB.ExecContext(ctx, query, args...)
		return opErr
	}

	err = RetryOperation(operation, DefaultMaxRetries)
	return result, err
}

func Exec(builder sq.Sqlizer) (sql.Result, error) {
	return ExecContext(context.Background(), builder)
}

func Query(builder sq.Sqlizer) (*sql.Rows, error) {
	return QueryContext(context.Background(), builder)
}

func QueryRow(builder sq.Sqlizer) *sql.Row {
	return QueryRowContext(context.Background(), builder)
}
