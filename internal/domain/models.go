package domain

import "time"

// Account represents a user's bank account.
type Account struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"` // Stored in cents to avoid floating point errors
	CreatedAt time.Time `json:"created_at"`
}

// Transfer represents a money transfer between two accounts.
type Transfer struct {
	ID            int       `json:"id"`
	FromAccountID int       `json:"from_account_id"`
	ToAccountID   int       `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}
