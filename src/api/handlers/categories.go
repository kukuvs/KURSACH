package handlers

import (
	"buch/src"
	"buch/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []src.Category
	// GORM: Найти все записи и отсортировать по имени
	if result := database.DB.Order("name").Find(&categories); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var cat src.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// GORM: Создать запись
	if result := database.DB.Create(&cat); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat src.Category

	// Сначала проверяем, существует ли категория
	if result := database.DB.First(&cat, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Читаем новые данные
	var input src.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// GORM: Обновить поля
	database.DB.Model(&cat).Updates(src.Category{
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Category updated", "category": cat})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// GORM: Удалить по ID
	// Unscoped() удаляет запись полностью из БД, а не делает soft delete
	// Если хочешь soft delete (чтобы можно было восстановить), убери .Unscoped()
	if result := database.DB.Unscoped().Delete(&src.Category{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
