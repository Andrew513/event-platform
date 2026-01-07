package postgres

import (
	"context"
	"time"
	"errors"
	"fmt"

	"github.com/Andrew513/event-platform/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type LedgerStore struct {
	pool *pgxpool.Pool
}

func NewLedgerStore(pool *pgxpool.Pool) *LedgerStore {
	return &LedgerStore{
		pool: pool,
	}
}

func (s *LedgerStore) ApplyEventTx(ctx context.Context, e domain.LedgerEvent) (alreadyProcessed bool, err error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO accounts (account_id, balance)
		VALUES ($1, 0)
		ON CONFLICT (account_id) DO NOTHING
	`, e.AccountID)
	if err != nil {
		return false, err
	}

	var delta float64
	switch e.Type {
	case "DEPOSIT":
		delta = e.Amount
	case "WITHDRAWAL":
		delta = -e.Amount
	default:
		return false, fmt.Errorf("unknown event type: %s", e.Type)
	}

	var insertedID string
	createdAt := e.Timestamp
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO ledger_entries (event_id, account_id, delta, type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (event_id) DO NOTHING
		RETURNING event_id
	`, e.EventID, e.AccountID, delta, e.Type, createdAt).Scan(&insertedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			alreadyProcessed = true
			err = tx.Commit(ctx)
			return alreadyProcessed, err
		}
		return false, err
	} 


	var newBalance float64
	err = tx.QueryRow(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE account_id = $2 AND (balance + $1) >= 0
		RETURNING balance
	`, delta, e.AccountID).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrInsufficientFunds
		}
		return false, err
	}

	err = tx.Commit(ctx)
	return false, err
}