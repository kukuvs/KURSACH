package api

import (
	"buch/src/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Статические файлы (CSS, JS)
	router.Static("/static", "./src/static")
	router.Static("/js", "./src/static/js")
	router.Static("/css", "./src/static/styles") // Обрати внимание: папка styles

	// HTML страницы
	router.GET("/", handlers.IndexPage)
	router.GET("/transactions", handlers.TransactionsPage)
	router.GET("/accounts", handlers.AccountsPage)
	router.GET("/reports", handlers.ReportsPage)
	router.GET("/categories", handlers.CategoriesPage)

	// API Группы
	api := router.Group("/api")
	{
		// Категории
		categories := api.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.POST("", handlers.CreateCategory)
			categories.PUT("/:id", handlers.UpdateCategory)
			categories.DELETE("/:id", handlers.DeleteCategory)
		}

		// Транзакции
		transactions := api.Group("/transactions")
		{
			transactions.GET("", handlers.GetTransactions)
			transactions.POST("", handlers.CreateTransaction)
			transactions.PUT("/:id", handlers.UpdateTransaction)
			transactions.DELETE("/:id", handlers.DeleteTransaction)
			transactions.GET("/by-date", handlers.GetTransactionsByDate)
		}

		// Счета
		accounts := api.Group("/accounts")
		{
			accounts.GET("", handlers.GetAccounts)
			accounts.POST("", handlers.CreateAccount)
			accounts.PUT("/:id", handlers.UpdateAccount)
			accounts.DELETE("/:id", handlers.DeleteAccount)
		}

		// Переводы
		transfers := api.Group("/transfers")
		{
			transfers.GET("", handlers.GetTransfers)
			transfers.POST("", handlers.CreateTransfer)
			transfers.DELETE("/:id", handlers.DeleteTransfer)
		}

		// Отчеты
		api.GET("/summary", handlers.GetSummary)
		api.GET("/stats/monthly", handlers.GetMonthlyStat)
	}

	return router
}
