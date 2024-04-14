package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a transactional function with the provided queries.
// It begins a new transaction using the given context and database connection.
// The provided function `fn` is called with a `Queries` instance that uses the transaction.
// If the function returns an error, the transaction is rolled back and the error is returned.
// Otherwise, the transaction is committed and nil is returned.
func (store *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	// Perform a transaction to create a transfer between two accounts and update their balances and transaction history.
	// The transfer amount is deducted from the 'FromAccountID' and added to the 'ToAccountID'.
	// Returns the created transfer, entries for both accounts, and updated balances for both accounts.
	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		// Create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create an entry for the 'FromAccountID' with a negative amount to deduct the transfer amount
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create an entry for the 'ToAccountID' with the transfer amount to add to the balance
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// Update the balances of the accounts
		// To prevent deadlock, we need to ensure that the accounts are updated in a consistent order.
		// One way to achieve this is by acquiring locks on the accounts in a specific order.
		// In this case, we can order the accounts based on their IDs and acquire locks in ascending order.
		// By doing so, we guarantee that no two transactions will acquire locks in opposite order, preventing deadlock.

		// Determine the order of updating the accounts
		if arg.FromAccountID < arg.ToAccountID {
			// If 'FromAccountID' is less than 'ToAccountID', deduct the transfer amount from 'FromAccountID' first
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)

			if err != nil {
				return err
			}

		} else {
			// If 'FromAccountID' is greater than or equal to 'ToAccountID', add the transfer amount to 'ToAccountID' first
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)

			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1, accountID2 int64,
	amount1, amount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		return
	}

	return
}
