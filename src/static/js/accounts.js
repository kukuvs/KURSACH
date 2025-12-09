const API_BASE = '/api';
let accountModal;

document.addEventListener('DOMContentLoaded', () => {
    accountModal = new bootstrap.Modal(document.getElementById('accountModal'));
    loadAccounts();
});

function openAccountModal() {
    document.getElementById('accountForm').reset();
    document.getElementById('accountId').value = '';
    document.getElementById('accountModalTitle').textContent = 'Новый счет';
    accountModal.show();
}

async function loadAccounts() {
    try {
        const response = await fetch(`${API_BASE}/accounts`);
        const accounts = await response.json();
        const container = document.getElementById('accountsContainer');
        container.innerHTML = '';

        if (accounts && accounts.length > 0) {
            accounts.forEach(acc => {
                const typeName = {
                    'cash': 'Наличные',
                    'bank': 'Банковский счет',
                    'savings': 'Сбережения'
                }[acc.account_type] || acc.account_type;

                const col = document.createElement('div');
                col.className = 'col-md-6 col-lg-4 mb-4';
                col.innerHTML = `
                    <div class="card border-0 shadow-sm h-100">
                        <div class="card-body">
                            <div class="d-flex justify-content-between align-items-start mb-2">
                                <h5 class="card-title mb-0">${acc.name}</h5>
                                <span class="badge bg-light text-dark border">${typeName}</span>
                            </div>
                            <h3 class="text-primary mb-3">${acc.balance.toFixed(2)} ₽</h3>
                        </div>
                        <div class="card-footer bg-white border-top-0 d-flex justify-content-end gap-2 pb-3">
                            <button class="btn btn-sm btn-outline-primary" onclick="editAccount(${acc.id})">
                                <i class="bi bi-pencil"></i>
                            </button>
                            <button class="btn btn-sm btn-outline-danger" onclick="deleteAccount(${acc.id})">
                                <i class="bi bi-trash"></i>
                            </button>
                        </div>
                    </div>
                `;
                container.appendChild(col);
            });
        } else {
            container.innerHTML = '<div class="col-12 text-center text-muted py-5">Нет счетов. Создайте первый!</div>';
        }
    } catch (error) {
        console.error('Error loading accounts:', error);
    }
}

async function saveAccount() {
    const id = document.getElementById('accountId').value;
    const data = {
        name: document.getElementById('accountName').value,
        account_type: document.getElementById('accountType').value,
        balance: parseFloat(document.getElementById('accountBalance').value)
    };

    try {
        let url = `${API_BASE}/accounts`;
        let method = 'POST';

        // Если ID есть - это редактирование
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
            accountModal.hide();
            loadAccounts();
        } else {
            alert('Ошибка при сохранении');
        }
    } catch (error) {
        console.error('Error saving account:', error);
    }
}

// ФУНКЦИЯ РЕДАКТИРОВАНИЯ
async function editAccount(id) {
    try {
        // Сначала получаем данные всех счетов (или конкретного, если API поддерживает GET /accounts/:id)
        // В нашем API пока GET /accounts возвращает все.
        const response = await fetch(`${API_BASE}/accounts`);
        const accounts = await response.json();
        
        // Ищем нужный счет (преобразуем id в строку для сравнения или используем ==)
        const acc = accounts.find(a => a.id == id);
        
        if (acc) {
            document.getElementById('accountId').value = acc.id;
            document.getElementById('accountName').value = acc.name;
            document.getElementById('accountType').value = acc.account_type;
            document.getElementById('accountBalance').value = acc.balance;
            
            document.getElementById('accountModalTitle').textContent = 'Редактирование счета';
            accountModal.show();
        } else {
            console.error('Account not found for id:', id);
        }
    } catch (e) {
        console.error('Error fetching account for edit:', e);
    }
}

async function deleteAccount(id) {
    if (!confirm('Удалить этот счет? Все связанные данные могут быть потеряны.')) return;
    
    try {
        const response = await fetch(`${API_BASE}/accounts/${id}`, { method: 'DELETE' });
        if (response.ok) loadAccounts();
        else alert('Ошибка удаления');
    } catch (e) {
        console.error(e);
    }
}
