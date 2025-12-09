package handlers

import (
	"buch/src"
	"buch/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSummary(c *gin.Context) {
	summary := src.Summary{
		ByCategory: make(map[string]float64),
	}

	// 1. Считаем общий доход
	// SELECT SUM(amount) FROM transactions WHERE transaction_type = 'income'
	var totalIncome float64
	database.DB.Model(&src.Transaction{}).
		Where("transaction_type = ?", "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)
	summary.TotalIncome = totalIncome

	// 2. Считаем общий расход
	var totalExpense float64
	database.DB.Model(&src.Transaction{}).
		Where("transaction_type = ?", "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)
	summary.TotalExpense = totalExpense

	// 3. Расходы по категориям
	// SELECT c.name, SUM(t.amount) ... GROUP BY c.name
	rows, err := database.DB.Table("transactions t").
		Select("c.name, SUM(t.amount)").
		Joins("JOIN categories c ON t.category_id = c.id").
		Where("t.transaction_type = ?", "expense"). // Обычно диаграммы строят по расходам
		Group("c.name").
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var amount float64
			if err := rows.Scan(&name, &amount); err == nil {
				summary.ByCategory[name] = amount
			}
		}
	}

	// 4. Счета и общий баланс по счетам
	if result := database.DB.Find(&summary.Accounts); result.Error == nil {
		for _, acc := range summary.Accounts {
			summary.TotalAccounts += acc.Balance
		}
	}

	// Баланс (Доходы - Расходы)
	summary.Balance = summary.TotalIncome - summary.TotalExpense

	c.JSON(http.StatusOK, summary)
}

func GetMonthlyStat(c *gin.Context) {
	// Статистика по месяцам для графика
	// Используем raw SQL, так как функции дат в разных БД отличаются (здесь под Postgres)

	type MonthlyResult struct {
		Month       string  `json:"month"`
		TransType   string  `json:"transaction_type"`
		TotalAmount float64 `json:"total"`
	}

	var results []MonthlyResult

	database.DB.Raw(`
		SELECT 
			TO_CHAR(transaction_date, 'YYYY-MM') as month,
			transaction_type as trans_type,
			SUM(amount) as total_amount
		FROM transactions
		GROUP BY 1, 2
		ORDER BY 1 DESC
		LIMIT 24
	`).Scan(&results)

	// Преобразуем в удобный формат для фронтенда
	type MonthlyStat struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
	}

	statsMap := make(map[string]*MonthlyStat)

	for _, r := range results {
		if _, exists := statsMap[r.Month]; !exists {
			statsMap[r.Month] = &MonthlyStat{Month: r.Month}
		}
		if r.TransType == "income" {
			statsMap[r.Month].Income = r.TotalAmount
		} else {
			statsMap[r.Month].Expense = r.TotalAmount
		}
	}

	c.JSON(http.StatusOK, statsMap)
}
