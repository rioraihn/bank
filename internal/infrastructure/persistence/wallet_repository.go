package persistence

import (
	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
	"context"
	"database/sql"
	"errors"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) GetWallet(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error) {
	query := `
		SELECT id, user_id, balance
		FROM wallets
		WHERE user_id = $1;
	`

	var walletID string
	var dbUserID string
	var balance int64

	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(
		&walletID,
		&dbUserID,
		&balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	// Convert to domain entities
	walletIDVO, err := valueobject.NewUserID(walletID)
	if err != nil {
		return nil, err
	}

	userIDVO, err := valueobject.NewUserID(dbUserID)
	if err != nil {
		return nil, err
	}

	balanceVO, err := valueobject.NewMoney(balance)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructWallet(walletIDVO, userIDVO, balanceVO), nil
}

func (r *WalletRepository) GetWalletForUpdate(ctx context.Context, tx *sql.Tx, userID valueobject.UserID) (*entity.Wallet, error) {
	query := `
		SELECT id, user_id, balance
		FROM wallets
		WHERE user_id = $1
		FOR UPDATE;
	`

	var walletID string
	var dbUserID string
	var balance int64

	err := tx.QueryRowContext(ctx, query, userID.String()).Scan(
		&walletID,
		&dbUserID,
		&balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	// Convert to domain entities
	walletIDVO, err := valueobject.NewUserID(walletID)
	if err != nil {
		return nil, err
	}

	userIDVO, err := valueobject.NewUserID(dbUserID)
	if err != nil {
		return nil, err
	}

	balanceVO, err := valueobject.NewMoney(balance)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructWallet(walletIDVO, userIDVO, balanceVO), nil
}

func (r *WalletRepository) UpdateWalletBalance(ctx context.Context, tx *sql.Tx, walletID valueobject.UserID, newBalance int64) error {
	query := `
		UPDATE wallets
		SET balance = $1, updated_at = NOW()
		WHERE id = $2;
	`

	_, err := tx.ExecContext(ctx, query, newBalance, walletID.String())
	return err
}
