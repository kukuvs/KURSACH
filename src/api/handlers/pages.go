package handlers

import (
	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	c.File("./src/static/index.html")
}

func TransactionsPage(c *gin.Context) {
	c.File("./src/static/transactions.html")
}

func AccountsPage(c *gin.Context) {
	c.File("./src/static/accounts.html")
}

func ReportsPage(c *gin.Context) {
	c.File("./src/static/reports.html")
}

func CategoriesPage(c *gin.Context) {
	c.File("./src/static/categories.html")
}
