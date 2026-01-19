package service

import (
	"context"

	"github.com/gabrielcau/ledger-system/internal/domain"
	"github.com/gabrielcau/ledger-system/internal/repository/postgres"
)

type LedgerService struct {
	repo *postgres.AccountRepository
}

func NewLedgerService(repo *postgres.AccountRepository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) Transfer(ctx context.Context, fromID, toID int, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}
	if fromID == toID {
		return domain.ErrSameAccount
	}

	transfer := domain.Transfer{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
	}

	return s.repo.TransferTx(ctx, transfer)
}

func (s *LedgerService) GetBalance(ctx context.Context, id int) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}
