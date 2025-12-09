package src

import (
	"time"

)

// Category - категория расходов/доходов
type Category struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"unique;not null" json:"name" binding:"required"`
	Description string    `json:"description"`
	Color       string    `gorm:"size:7" json:"color"`

	// Связь: одна категория может иметь много транзакций
	Transactions []Transaction `json:"-"`
}



type Transaction struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	CategoryID      uint      `json:"category_id" binding:"required"`
	Category        Category  `gorm:"foreignKey:CategoryID" json:"category" binding:"-"`

	// Добавили связь со счетом
	AccountID       uint      `json:"account_id" binding:"required"`
	Account         Account   `gorm:"foreignKey:AccountID" json:"account" binding:"-"`

	Description     string    `json:"description"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount" binding:"required,gt=0"`
	TransactionType string    `gorm:"size:20;not null" json:"transaction_type" binding:"required,oneof=income expense"`
	TransactionDate string    `gorm:"type:date;not null" json:"transaction_date" binding:"required"`
	Notes           string    `json:"notes"`
}

// Account - счет
type Account struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string  `gorm:"not null" json:"name" binding:"required"`
	Balance     float64 `gorm:"type:decimal(12,2);default:0" json:"balance"`
	AccountType string  `gorm:"size:50;not null" json:"account_type" binding:"required,oneof=cash bank savings"`
}

// Transfer - перевод между счетами
type Transfer struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	FromAccountID uint    `json:"from_account_id" binding:"required"`
	FromAccount   Account `gorm:"foreignKey:FromAccountID" json:"from_account"`

	ToAccountID uint    `json:"to_account_id" binding:"required"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID" json:"to_account"`

	Amount       float64   `gorm:"type:decimal(10,2);not null" json:"amount" binding:"required,gt=0"`
	Description  string    `json:"description"`
	TransferDate time.Time `gorm:"type:date;not null" json:"transfer_date" binding:"required"`
}

// Summary - сводка (не таблица БД, просто структура для JSON)
type Summary struct {
	TotalIncome   float64            `json:"total_income"`
	TotalExpense  float64            `json:"total_expense"`
	Balance       float64            `json:"balance"`
	ByCategory    map[string]float64 `json:"by_category"`
	Accounts      []Account          `json:"accounts"`
	TotalAccounts float64            `json:"total_accounts"`
}
