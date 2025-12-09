package handlers

import (
	"buch/src"
	"buch/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTransfers(c *gin.Context) {
	var transfers []src.Transfer
	// Загружаем связанные счета (от кого и кому) для отображения имен
	if result := database.DB.
		Preload("FromAccount").
		Preload("ToAccount").
		Order("transfer_date DESC").
		Find(&transfers); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transfers)
}

func CreateTransfer(c *gin.Context) {
	var tr src.Transfer
	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Начинаем транзакцию базы данных
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Проверяем баланс отправителя
		var fromAcc src.Account
		if err := tx.First(&fromAcc, tr.FromAccountID).Error; err != nil {
			return err
		}

		if fromAcc.Balance < tr.Amount {
			// Можно вернуть кастомную ошибку, но для простоты вернем стандартную
			return gorm.ErrInvalidData // "Недостаточно средств" можно обработать на клиенте
		}

		// 2. Списываем средства
		if err := tx.Model(&src.Account{}).Where("id = ?", tr.FromAccountID).
			Update("balance", gorm.Expr("balance - ?", tr.Amount)).Error; err != nil {
			return err
		}

		// 3. Зачисляем средства
		if err := tx.Model(&src.Account{}).Where("id = ?", tr.ToAccountID).
			Update("balance", gorm.Expr("balance + ?", tr.Amount)).Error; err != nil {
			return err
		}

		// 4. Создаем запись о переводе
		if err := tx.Create(&tr).Error; err != nil {
			return err
		}

		return nil // Commit transaction
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transfer failed: " + err.Error()})
		return
	}

	// Загружаем данные для ответа
	database.DB.Preload("FromAccount").Preload("ToAccount").First(&tr, tr.ID)
	c.JSON(http.StatusCreated, tr)
}

func DeleteTransfer(c *gin.Context) {
	id := c.Param("id")

	// При удалении перевода по-хорошему надо бы вернуть деньги обратно,
	// но для простоты пока просто удаляем запись.
	// Если нужно вернуть деньги, это тоже надо делать в транзакции.

	if result := database.DB.Unscoped().Delete(&src.Transfer{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer deleted"})
}
