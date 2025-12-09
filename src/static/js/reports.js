const API_BASE = "/api";
let incomeExpenseChart = null;
let categoryChart = null;
let monthlyChart = null;

document.addEventListener("DOMContentLoaded", () => {
	loadReports();
});

async function loadReports() {
	try {
		const summaryRes = await fetch(`${API_BASE}/summary`);
		const summary = await summaryRes.json();

		renderIncomeExpenseChart(summary.total_income, summary.total_expense);
		renderCategoryChart(summary.by_category);

		const statsRes = await fetch(`${API_BASE}/stats/monthly`);
		const stats = await statsRes.json();

		renderMonthlyChart(stats);
	} catch (error) {
		console.error("Error loading reports:", error);
	}
}

function renderIncomeExpenseChart(income, expense) {
	const ctx = document.getElementById("incomeExpenseChart").getContext("2d");
	if (incomeExpenseChart) incomeExpenseChart.destroy();

	incomeExpenseChart = new Chart(ctx, {
		type: "doughnut",
		data: {
			labels: ["Доходы", "Расходы"],
			datasets: [
				{
					data: [income, expense],
					backgroundColor: ["#198754", "#dc3545"],
					hoverOffset: 4,
				},
			],
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
		},
	});
}

function renderCategoryChart(byCategory) {
	const ctx = document.getElementById("categoryChart").getContext("2d");
	if (categoryChart) categoryChart.destroy();

	// Защита от null
	if (!byCategory) byCategory = {};

	const labels = Object.keys(byCategory);
	const data = Object.values(byCategory);

	categoryChart = new Chart(ctx, {
		type: "bar",
		data: {
			labels: labels,
			datasets: [
				{
					label: "Расходы",
					data: data,
					backgroundColor: "#0d6efd",
					borderRadius: 5,
				},
			],
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			indexAxis: "y",
			scales: {
				x: { beginAtZero: true },
			},
		},
	});
}

function renderMonthlyChart(statsMap) {
	const ctx = document.getElementById("monthlyChart").getContext("2d");
	if (monthlyChart) monthlyChart.destroy();

	if (!statsMap) statsMap = {};

	const months = Object.keys(statsMap).sort();

	const incomeData = months.map((m) => statsMap[m].income || 0);
	const expenseData = months.map((m) => statsMap[m].expense || 0);

	monthlyChart = new Chart(ctx, {
		type: "line",
		data: {
			labels: months,
			datasets: [
				{
					label: "Доходы",
					data: incomeData,
					borderColor: "#198754",
					backgroundColor: "rgba(25, 135, 84, 0.1)",
					fill: true,
					tension: 0.3,
				},
				{
					label: "Расходы",
					data: expenseData,
					borderColor: "#dc3545",
					backgroundColor: "rgba(220, 53, 69, 0.1)",
					fill: true,
					tension: 0.3,
				},
			],
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			scales: {
				y: { beginAtZero: true },
			},
		},
	});
}
