package handlers

import (
	"buch/src"
	"buch/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	var transactions []src.Transaction
	// Добавили Preload("Account")
	if result := database.DB.Preload("Category").Preload("Account").Order("transaction_date DESC").Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func CreateTransaction(c *gin.Context) {
	var trans src.Transaction

	// 1. Валидация входных данных
	if err := c.ShouldBindJSON(&trans); err != nil {
		// Логируем ошибку для отладки
		// fmt.Println("Validation error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Начинаем транзакцию БД (чтобы если создание упадет, баланс не обновился, и наоборот)
	tx := database.DB.Begin()

	// 3. Проверяем существование счета
	var account src.Account
	if err := tx.First(&account, trans.AccountID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found"})
		return
	}

	// 4. Создаем запись транзакции
	if err := tx.Create(&trans).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Обновляем баланс счета
	switch trans.TransactionType {
	case "income":
		account.Balance += trans.Amount
	case "expense":
		account.Balance -= trans.Amount
	}

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
		return
	}

	// 6. Фиксируем транзакцию БД
	tx.Commit()

	// 7. Загружаем связанные данные (Категорию и Счет) для красивого ответа фронтенду
	database.DB.Preload("Category").Preload("Account").First(&trans, trans.ID)

	c.JSON(http.StatusCreated, trans)
}

func UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	var trans src.Transaction

	if result := database.DB.First(&trans, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input src.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем поля
	if result := database.DB.Model(&trans).Updates(src.Transaction{
		CategoryID:      input.CategoryID,
		Description:     input.Description,
		Amount:          input.Amount,
		TransactionType: input.TransactionType,
		TransactionDate: input.TransactionDate,
		Notes:           input.Notes,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated", "transaction": trans})
}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	// Unscoped() для полного удаления
	if result := database.DB.Unscoped().Delete(&src.Transaction{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

func GetTransactionsByDate(c *gin.Context) {
	startDate := c.Query("start")
	endDate := c.Query("end")

	var transactions []src.Transaction

	// Фильтрация по дате
	if result := database.DB.Preload("Category").Preload("Account"). // Добавили здесь
										Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
										Order("transaction_date DESC").
										Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
