const API_BASE = '/api';
let transactionModal;

document.addEventListener('DOMContentLoaded', () => {
    transactionModal = new bootstrap.Modal(document.getElementById('transactionModal'));
    
    // Устанавливаем текущую дату
    document.getElementById('transDate').valueAsDate = new Date();
    
    // Сначала загружаем справочники (категории и счета)
    Promise.all([loadCategories(), loadAccountsForSelect()])
        .then(() => {
            // И только потом загружаем список транзакций
            loadTransactions();
        });

    // Вешаем обработчики на фильтры
    document.getElementById('filterType').addEventListener('change', loadTransactions);
    document.getElementById('filterCategory').addEventListener('change', loadTransactions);
    document.getElementById('filterStartDate').addEventListener('change', loadTransactions);
    document.getElementById('filterEndDate').addEventListener('change', loadTransactions);
});

// Загрузка категорий в оба селекта (в фильтр и в модалку)
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/categories`);
        if (!response.ok) throw new Error('Failed to fetch categories');
        
        const categories = await response.json();

        // Селекты, куда нужно добавить категории
        const selects = [
            { el: document.getElementById('transCategory'), defaultText: 'Выберите категорию' },
            { el: document.getElementById('filterCategory'), defaultText: 'Все категории' }
        ];
        
        selects.forEach(item => {
            const select = item.el;
            if (!select) return;

            // Очищаем и добавляем дефолтную опцию
            select.innerHTML = `<option value="" disabled selected>${item.defaultText}</option>`;
            
            // Если это фильтр, то "Все категории" должна быть выбираемой (value="")
            if (select.id === 'filterCategory') {
                 select.innerHTML = `<option value="">Все категории</option>`;
            }

            if (categories && categories.length > 0) {
                categories.forEach(cat => {
                    const option = document.createElement('option');
                    option.value = cat.id;
                    // Добавляем цветной кружок (символически текстом, т.к. в select нельзя HTML)
                    option.textContent = cat.name; 
                    select.appendChild(option);
                });
            }
        });
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

// Загрузка счетов
async function loadAccountsForSelect() {
    try {
        const response = await fetch(`${API_BASE}/accounts`);
        if (!response.ok) throw new Error('Failed to fetch accounts');
        
        const accounts = await response.json();
        const select = document.getElementById('transAccount');
        
        if (!select) return;

        select.innerHTML = '<option value="" disabled selected>Выберите счет</option>';

        if (accounts && accounts.length > 0) {
            accounts.forEach(acc => {
                const option = document.createElement('option');
                option.value = acc.id;
                option.textContent = `${acc.name} (${acc.balance.toFixed(2)} ₽)`;
                select.appendChild(option);
            });
        }
    } catch (error) {
        console.error('Error loading accounts:', error);
    }
}

// Загрузка списка транзакций
async function loadTransactions() {
    try {
        let url = `${API_BASE}/transactions`;
        const startDate = document.getElementById('filterStartDate').value;
        const endDate = document.getElementById('filterEndDate').value;

        if (startDate && endDate) {
            url = `${API_BASE}/transactions/by-date?start=${startDate}&end=${endDate}`;
        }

        const response = await fetch(url);
        if (!response.ok) throw new Error('Failed to fetch transactions');
        
        const transactions = await response.json();

        const tbody = document.getElementById('transactionsTableBody');
        tbody.innerHTML = '';

        if (transactions && transactions.length > 0) {
            let filtered = transactions;

            // Фильтрация на клиенте (дополнительно к серверной по дате)
            const typeFilter = document.getElementById('filterType').value;
            if (typeFilter) {
                filtered = filtered.filter(t => t.transaction_type === typeFilter);
            }

            const catFilter = document.getElementById('filterCategory').value;
            if (catFilter) {
                filtered = filtered.filter(t => t.category_id == catFilter);
            }

            filtered.forEach(trans => {
                const row = document.createElement('tr');
                const isIncome = trans.transaction_type === 'income';
                const amountClass = isIncome ? 'text-success' : 'text-danger';
                const sign = isIncome ? '+' : '-';
                const typeName = isIncome ? 'Доход' : 'Расход';
                const categoryName = trans.category ? trans.category.name : 'Без категории';
                const accountName = trans.account ? trans.account.name : '-';

                // Форматируем дату красиво
                const dateObj = new Date(trans.transaction_date);
                const dateStr = dateObj.toLocaleDateString('ru-RU');

                row.innerHTML = `
                    <td>${dateStr}</td>
                    <td><span class="badge bg-secondary">${categoryName}</span></td>
                    <td><small class="text-muted">${accountName}</small></td>
                    <td>${trans.description || '-'}</td>
                    <td>${typeName}</td>
                    <td class="text-end ${amountClass} fw-bold">${sign}${trans.amount.toFixed(2)} ₽</td>
                    <td class="text-end">
                        <button class="btn btn-sm btn-outline-danger" onclick="deleteTransaction(${trans.id})">
                            <i class="bi bi-trash"></i>
                        </button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="7" class="text-center text-muted">Нет данных</td></tr>';
        }
    } catch (error) {
        console.error('Error loading transactions:', error);
    }
}

function openTransactionModal() {
    document.getElementById('transactionForm').reset();
    document.getElementById('transId').value = '';
    document.getElementById('modalTitle').textContent = 'Новая транзакция';
    document.getElementById('transDate').valueAsDate = new Date();
    transactionModal.show();
}

async function saveTransaction() {
    const id = document.getElementById('transId').value;
    const dateStr = document.getElementById('transDate').value;

    const data = {
        transaction_type: document.getElementById('transType').value,
        category_id: parseInt(document.getElementById('transCategory').value),
        account_id: parseInt(document.getElementById('transAccount').value),
        transaction_date: dateStr,
        amount: parseFloat(document.getElementById('transAmount').value),
        description: document.getElementById('transDescription').value,
        notes: document.getElementById('transNotes').value
    };

    if (isNaN(data.category_id)) {
        alert('Выберите категорию!');
        return;
    }
    if (isNaN(data.account_id)) {
        alert('Выберите счет!');
        return;
    }
    if (isNaN(data.amount) || data.amount <= 0) {
        alert('Введите корректную сумму!');
        return;
    }

    try {
        let url = `${API_BASE}/transactions`;
        let method = 'POST';

        if (id) {
            url += '/' + id;
            method = 'PUT';
        }

        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            transactionModal.hide();
            loadTransactions();
            // Можно перезагрузить счета, чтобы обновить балансы в списке (если нужно)
            loadAccountsForSelect();
        } else {
            const err = await response.json();
            alert('Ошибка: ' + (err.error || 'Неизвестная ошибка'));
        }
    } catch (error) {
        console.error('Error saving transaction:', error);
    }
}

async function deleteTransaction(id) {
    if (!confirm('Вы уверены? Удаление транзакции не вернет деньги на счет (пока не реализовано).')) return;

    try {
        const response = await fetch(`${API_BASE}/transactions/${id}`, { method: 'DELETE' });
        if (response.ok) {
            loadTransactions();
        } else {
            alert('Ошибка удаления');
        }
    } catch (error) {
        console.error('Error deleting:', error);
    }
}
