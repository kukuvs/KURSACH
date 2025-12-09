package handlers

import (
	"buch/src"
	"buch/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccounts(c *gin.Context) {
	var accounts []src.Account
	// Получить все счета, отсортированные по имени
	if result := database.DB.Order("name").Find(&accounts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func CreateAccount(c *gin.Context) {
	var acc src.Account
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&acc); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	var acc src.Account

	if result := database.DB.First(&acc, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var input src.Account
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Model(&acc).Updates(src.Account{
		Name:        input.Name,
		Balance:     input.Balance,
		AccountType: input.AccountType,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account updated", "account": acc})
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("id")

	if result := database.DB.Unscoped().Delete(&src.Account{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}
