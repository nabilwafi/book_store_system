package repository

import (
	"github.com/nabil/book-store-system/pkg/logger"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}

// BeginTransaction starts a new database transaction
func (r *transactionRepositoryImpl) BeginTransaction() (*gorm.DB, error) {
	logger.Infof("Starting new database transaction")
	tx := r.db.Begin()
	if tx.Error != nil {
		logger.Errorf("Failed to begin transaction: %v", tx.Error)
		return nil, tx.Error
	}
	logger.Infof("Database transaction started successfully")
	return tx, nil
}

// CommitTransaction commits the transaction
func (r *transactionRepositoryImpl) CommitTransaction(tx *gorm.DB) error {
	logger.Infof("Committing database transaction")
	err := tx.Commit().Error
	if err != nil {
		logger.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	logger.Infof("Database transaction committed successfully")
	return nil
}

// RollbackTransaction rolls back the transaction
func (r *transactionRepositoryImpl) RollbackTransaction(tx *gorm.DB) error {
	logger.Infof("Rolling back database transaction")
	err := tx.Rollback().Error
	if err != nil {
		logger.Errorf("Failed to rollback transaction: %v", err)
		return err
	}
	logger.Infof("Database transaction rolled back successfully")
	return nil
}

// WithTransaction executes a function within a transaction
func (r *transactionRepositoryImpl) WithTransaction(fn func(tx *gorm.DB) error) error {
	logger.Infof("Executing function within transaction")

	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("Panic occurred, rolling back transaction: %v", p)
			r.RollbackTransaction(tx)
			panic(p) // re-throw panic after rollback
		}
	}()

	err = fn(tx)
	if err != nil {
		logger.Errorf("Function execution failed, rolling back transaction: %v", err)
		if rollbackErr := r.RollbackTransaction(tx); rollbackErr != nil {
			logger.Errorf("Failed to rollback after error: %v", rollbackErr)
		}
		return err
	}

	return r.CommitTransaction(tx)
}
