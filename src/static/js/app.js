const API_BASE = "/api";

// Загрузка сводной информации (доходы, расходы, баланс)
async function loadSummary() {
	try {
		const response = await fetch(`${API_BASE}/summary`);
		if (!response.ok) throw new Error("Network response was not ok");

		const summary = await response.json();

		// Обновляем карточки статистики
		document.getElementById("totalIncome").textContent =
			summary.total_income.toFixed(2) + " ₽";
		document.getElementById("totalExpense").textContent =
			summary.total_expense.toFixed(2) + " ₽";

		const balanceEl = document.getElementById("balance");
		balanceEl.textContent = summary.balance.toFixed(2) + " ₽";

		// Меняем цвет баланса в зависимости от значения
		if (summary.balance < 0) {
			balanceEl.parentElement.parentElement.classList.remove(
				"bg-primary"
			);
			balanceEl.parentElement.parentElement.classList.add("bg-danger");
		} else {
			balanceEl.parentElement.parentElement.classList.remove("bg-danger");
			balanceEl.parentElement.parentElement.classList.add("bg-primary");
		}

		// Заполняем список счетов справа
		const accountsList = document.getElementById("accountsList");
		accountsList.innerHTML = "";

		if (summary.accounts && summary.accounts.length > 0) {
			summary.accounts.forEach((acc) => {
				const item = document.createElement("div");
				item.className =
					"list-group-item d-flex justify-content-between align-items-center";
				item.innerHTML = `
                    <div>
                        <div class="fw-bold">${acc.name}</div>
                        <small class="text-muted">${getAccountTypeName(
							acc.account_type
						)}</small>
                    </div>
                    <span class="badge bg-primary rounded-pill">${acc.balance.toFixed(
						2
					)} ₽</span>
                `;
				accountsList.appendChild(item);
			});
		} else {
			accountsList.innerHTML =
				'<div class="text-center p-3 text-muted">Нет счетов</div>';
		}
	} catch (error) {
		console.error("Error loading summary:", error);
	}
}

// Загрузка последних транзакций
async function loadRecentTransactions() {
	try {
		const response = await fetch(`${API_BASE}/transactions`);
		if (!response.ok) throw new Error("Network response was not ok");

		const transactions = await response.json();
		const tbody = document.getElementById("transactionsBody");
		tbody.innerHTML = "";

		if (transactions && transactions.length > 0) {
			// Берем только 5 последних
			transactions.slice(0, 5).forEach((trans) => {
				const row = document.createElement("tr");

				const isIncome = trans.transaction_type === "income";
				const amountClass = isIncome
					? "amount-income"
					: "amount-expense";
				const sign = isIncome ? "+" : "-";
				const categoryName = trans.category
					? trans.category.name
					: "Без категории";

				row.innerHTML = `
                    <td>${formatDate(trans.transaction_date)}</td>
                    <td><span class="badge bg-secondary">${categoryName}</span></td>
                    <td>${trans.description || "-"}</td>
                    <td class="text-end ${amountClass}">${sign}${trans.amount.toFixed(
					2
				)} ₽</td>
                `;
				tbody.appendChild(row);
			});
		} else {
			tbody.innerHTML =
				'<tr><td colspan="4" class="text-center text-muted">Нет транзакций</td></tr>';
		}
	} catch (error) {
		console.error("Error loading transactions:", error);
	}
}

function formatDate(dateStr) {
	if (!dateStr) return "";
	const date = new Date(dateStr);
	return date.toLocaleDateString("ru-RU");
}

function getAccountTypeName(type) {
	const types = {
		cash: "Наличные",
		bank: "Банковский счет",
		savings: "Сбережения",
	};
	return types[type] || type;
}

document.addEventListener("DOMContentLoaded", () => {
	loadSummary();
	loadRecentTransactions();
});
