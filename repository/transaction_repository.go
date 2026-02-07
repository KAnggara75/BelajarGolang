package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrProductNotFoundInTx = errors.New("product not found")
	ErrEmptyCheckout       = errors.New("checkout items cannot be empty")
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Checkout(ctx context.Context, req models.CheckoutRequest) (models.Transaction, error)
	GetAll(ctx context.Context) ([]models.Transaction, error)
	GetByID(ctx context.Context, id int) (models.Transaction, error)
}

// transactionRepository implements TransactionRepository using PostgreSQL
type transactionRepository struct {
	db *pgx.Conn
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *pgx.Conn) TransactionRepository {
	return &transactionRepository{db: db}
}

// Checkout processes a checkout request with proper transaction handling
func (r *transactionRepository) Checkout(ctx context.Context, req models.CheckoutRequest) (models.Transaction, error) {
	if len(req.Items) == 0 {
		return models.Transaction{}, ErrEmptyCheckout
	}

	// Start a database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback if not committed

	var totalAmount int
	var details []models.TransactionDetail

	// Process each item
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return models.Transaction{}, ErrInvalidQuantity
		}

		// Lock the product row for update to prevent race conditions
		var productID int
		var productName string
		var price float64
		var stock int

		query := `SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE`
		err := tx.QueryRow(ctx, query, item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.Transaction{}, ErrProductNotFoundInTx
			}
			return models.Transaction{}, fmt.Errorf("failed to get product: %w", err)
		}

		// Check if sufficient stock
		if stock < item.Quantity {
			return models.Transaction{}, fmt.Errorf("%w: product %s has only %d in stock, requested %d",
				ErrInsufficientStock, productName, stock, item.Quantity)
		}

		// Calculate subtotal (convert price to int cents)
		subtotal := int(price*100) * item.Quantity
		totalAmount += subtotal

		// Update product stock
		updateQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateQuery, item.Quantity, item.ProductID)
		if err != nil {
			return models.Transaction{}, fmt.Errorf("failed to update stock: %w", err)
		}

		// Store detail for later insertion
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction record
	insertTxQuery := `INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at`
	var transaction models.Transaction
	err = tx.QueryRow(ctx, insertTxQuery, totalAmount).Scan(&transaction.ID, &transaction.CreatedAt)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	transaction.TotalAmount = totalAmount

	// Insert transaction details
	for i := range details {
		insertDetailQuery := `INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRow(ctx, insertDetailQuery,
			transaction.ID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return models.Transaction{}, fmt.Errorf("failed to create transaction detail: %w", err)
		}
		details[i].TransactionID = transaction.ID
	}

	transaction.Details = details

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return models.Transaction{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// GetAll returns all transactions with their details
func (r *transactionRepository) GetAll(ctx context.Context) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.total_amount, t.created_at,
		       td.id, td.product_id, p.name, td.quantity, td.subtotal
		FROM transactions t
		LEFT JOIN transaction_details td ON t.id = td.transaction_id
		LEFT JOIN products p ON td.product_id = p.id
		ORDER BY t.created_at DESC, td.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactionMap := make(map[int]*models.Transaction)
	var transactionIDs []int

	for rows.Next() {
		var txID int
		var totalAmount int
		var createdAt interface{}
		var detailID, productID, quantity, subtotal *int
		var productName *string

		err := rows.Scan(&txID, &totalAmount, &createdAt, &detailID, &productID, &productName, &quantity, &subtotal)
		if err != nil {
			return nil, err
		}

		// Get or create transaction
		tx, exists := transactionMap[txID]
		if !exists {
			tx = &models.Transaction{
				ID:          txID,
				TotalAmount: totalAmount,
				Details:     []models.TransactionDetail{},
			}
			// Handle created_at timestamp
			if t, ok := createdAt.(interface{ Time() interface{} }); ok {
				if timeVal, ok := t.Time().(interface{ String() string }); ok {
					_ = timeVal // Use if needed
				}
			}
			transactionMap[txID] = tx
			transactionIDs = append(transactionIDs, txID)
		}

		// Add detail if exists
		if detailID != nil && productID != nil {
			detail := models.TransactionDetail{
				ID:            *detailID,
				TransactionID: txID,
				ProductID:     *productID,
				Quantity:      *quantity,
				Subtotal:      *subtotal,
			}
			if productName != nil {
				detail.ProductName = *productName
			}
			tx.Details = append(tx.Details, detail)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice maintaining order
	var transactions []models.Transaction
	for _, id := range transactionIDs {
		transactions = append(transactions, *transactionMap[id])
	}

	return transactions, nil
}

// GetByID returns a transaction by ID with its details
func (r *transactionRepository) GetByID(ctx context.Context, id int) (models.Transaction, error) {
	// Get transaction
	txQuery := `SELECT id, total_amount, created_at FROM transactions WHERE id = $1`
	var transaction models.Transaction
	err := r.db.QueryRow(ctx, txQuery, id).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Transaction{}, ErrTransactionNotFound
		}
		return models.Transaction{}, err
	}

	// Get transaction details
	detailsQuery := `
		SELECT td.id, td.transaction_id, td.product_id, p.name, td.quantity, td.subtotal
		FROM transaction_details td
		LEFT JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id = $1
		ORDER BY td.id
	`

	rows, err := r.db.Query(ctx, detailsQuery, id)
	if err != nil {
		return models.Transaction{}, err
	}
	defer rows.Close()

	var details []models.TransactionDetail
	for rows.Next() {
		var detail models.TransactionDetail
		var productName *string
		err := rows.Scan(&detail.ID, &detail.TransactionID, &detail.ProductID, &productName, &detail.Quantity, &detail.Subtotal)
		if err != nil {
			return models.Transaction{}, err
		}
		if productName != nil {
			detail.ProductName = *productName
		}
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return models.Transaction{}, err
	}

	transaction.Details = details
	return transaction, nil
}
