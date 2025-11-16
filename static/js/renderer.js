const API_BASE = '/api';

let currentUser = null;

// Управление табами авторизации
function showTab(tabName) {
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
    
    event.target.classList.add('active');
    document.getElementById(`${tabName}-tab`).classList.add('active');
}

// Регистрация
async function register() {
    const cipher = document.getElementById('reg-cipher').value;
    const password = document.getElementById('reg-password').value;
    const name = document.getElementById('reg-name').value;
    const group = document.getElementById('reg-group').value;

    if (!cipher || !password || !name || !group) {
        showMessage('Заполните все поля', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/users/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ cipher, password, name, group })
        });

        const result = await response.json();
        
        if (result.success) {
            showMessage(result.message, 'success');
            showTab('login');
            // Очищаем поля
            document.getElementById('reg-cipher').value = '';
            document.getElementById('reg-password').value = '';
            document.getElementById('reg-name').value = '';
            document.getElementById('reg-group').value = '';
        } else {
            showMessage(result.error, 'error');
        }
    } catch (error) {
        showMessage('Ошибка регистрации: ' + error.message, 'error');
    }
}

// Вход
async function login() {
    const cipher = document.getElementById('login-cipher').value;
    const password = document.getElementById('login-password').value;

    if (!cipher || !password) {
        showMessage('Введите шифр и пароль', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/users/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ cipher, password })
        });

        const result = await response.json();
        
        if (result.success) {
            currentUser = result.user;
            showMainSection();
            loadUserData();
            showMessage('Успешный вход!', 'success');
        } else {
            showMessage(result.error, 'error');
        }
    } catch (error) {
        showMessage('Ошибка входа: ' + error.message, 'error');
    }
}

// Показ главной секции
function showMainSection() {
    document.getElementById('auth-section').classList.remove('active');
    document.getElementById('main-section').classList.add('active');
}

// Загрузка данных пользователя
async function loadUserData() {
    if (!currentUser) return;

    try {
        const response = await fetch(`${API_BASE}/users/${currentUser.cipher}`);
        const data = await response.json();
        
        if (response.ok) {
            document.getElementById('balance-amount').textContent = 
                `${data.user.balance} Кибиков`;
            document.getElementById('user-name').textContent = data.user.name;
            
            updateTransactionsList(data.transactions);
            updateRecipientsList();
            loadTasks();
        } else {
            showMessage(data.error, 'error');
        }
    } catch (error) {
        showMessage('Ошибка загрузки данных: ' + error.message, 'error');
    }
}

// Обновление списка транзакций
function updateTransactionsList(transactions) {
    const container = document.getElementById('transactions-list');
    
    if (!transactions || transactions.length === 0) {
        container.innerHTML = `
            <div class="transaction-item">
                <div class="transaction-info">
                    <div class="transaction-name">Операций пока нет</div>
                    <div class="transaction-date">Совершите первую операцию</div>
                </div>
            </div>
        `;
        return;
    }
    
    container.innerHTML = transactions.map(transaction => `
        <div class="transaction-item">
            <div class="transaction-info">
                <div class="transaction-name">${transaction.recipient_name}</div>
                <div class="transaction-date">${new Date(transaction.date).toLocaleDateString('ru-RU')}</div>
            </div>
            <div class="transaction-amount ${transaction.to_cipher === currentUser.cipher ? 'positive' : 'negative'}">
                ${transaction.to_cipher === currentUser.cipher ? '+' : '-'}${transaction.amount} К
            </div>
        </div>
    `).join('');
}

// Обновление списка получателей
async function updateRecipientsList() {
    const select = document.getElementById('transfer-recipient');
    
    try {
        const response = await fetch(`${API_BASE}/users`);
        const users = await response.json();
        
        select.innerHTML = '<option value="">Выберите получателя</option>';
        
        // Фильтруем текущего пользователя
        const otherUsers = users.filter(user => user.cipher !== currentUser.cipher);
        
        otherUsers.forEach(user => {
            const option = document.createElement('option');
            option.value = user.cipher;
            option.textContent = `${user.name} (${user.group_name})`;
            option.dataset.group = user.group_name;
            select.appendChild(option);
        });
    } catch (error) {
        showMessage('Ошибка загрузки пользователей: ' + error.message, 'error');
    }
}

// Фильтрация получателей по группе
function filterRecipients() {
    const group = document.getElementById('transfer-group').value;
    const select = document.getElementById('transfer-recipient');
    
    Array.from(select.options).forEach(option => {
        if (option.value === '') return;
        option.style.display = !group || option.dataset.group === group ? 'block' : 'none';
    });
}

// Выполнение перевода
async function makeTransfer() {
    if (!currentUser) return;

    const toCipher = document.getElementById('transfer-recipient').value;
    const amount = parseFloat(document.getElementById('transfer-amount').value);
    const recipientOption = document.getElementById('transfer-recipient').selectedOptions[0];

    if (!toCipher || !amount || amount <= 0) {
        showMessage('Заполните все поля корректно', 'error');
        return;
    }

    if (amount > currentUser.balance) {
        showMessage('Недостаточно средств', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/transfers/${currentUser.cipher}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                to_cipher: toCipher, 
                amount: amount 
            })
        });

        const result = await response.json();
        
        if (result.success) {
            showMessage(result.message, 'success');
            document.getElementById('transfer-amount').value = '';
            // Обновляем баланс текущего пользователя
            currentUser.balance = result.new_balance;
            loadUserData();
        } else {
            showMessage(result.error, 'error');
        }
    } catch (error) {
        showMessage('Ошибка перевода: ' + error.message, 'error');
    }
}

// Загрузка заданий
async function loadTasks() {
    try {
        const response = await fetch(`${API_BASE}/tasks`);
        const tasks = await response.json();
        
        const container = document.getElementById('tasks-container');
        
        container.innerHTML = tasks.map(task => `
            <div class="task-item">
                <div class="task-info">
                    <div class="task-title">${task.title}</div>
                    <div class="task-status">${task.completed ? '✅ Выполнено' : '⏳ Доступно'}</div>
                </div>
                <div class="task-reward">+${task.reward} К</div>
                ${!task.completed ? `<button onclick="completeTask(${task.id})" class="btn-primary" style="padding: 5px 10px; font-size: 12px;">Выполнить</button>` : ''}
            </div>
        `).join('');
    } catch (error) {
        showMessage('Ошибка загрузки заданий: ' + error.message, 'error');
    }
}

// Выполнение задания
async function completeTask(taskId) {
    if (!currentUser) return;

    try {
        const response = await fetch(`${API_BASE}/tasks/${currentUser.cipher}/complete/${taskId}`, {
            method: 'POST'
        });

        const result = await response.json();
        
        if (result.success) {
            showMessage(result.message, 'success');
            // Обновляем баланс
            currentUser.balance = result.new_balance;
            loadUserData();
            loadTasks();
        } else {
            showMessage(result.error, 'error');
        }
    } catch (error) {
        showMessage('Ошибка выполнения задания: ' + error.message, 'error');
    }
}

// Управление вкладками главного меню
function showMainTab(tabName) {
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.main-tab').forEach(tab => tab.classList.remove('active'));
    
    event.target.classList.add('active');
    document.getElementById(`${tabName}-tab`).classList.add('active');
}

// Выход
function logout() {
    currentUser = null;
    document.getElementById('main-section').classList.remove('active');
    document.getElementById('auth-section').classList.add('active');
    
    document.getElementById('login-cipher').value = '';
    document.getElementById('login-password').value = '';
}

// Вспомогательная функция для показа сообщений
function showMessage(text, type) {
    const existingMessage = document.querySelector('.message');
    if (existingMessage) {
        existingMessage.remove();
    }
    
    const message = document.createElement('div');
    message.className = `message ${type}`;
    message.textContent = text;
    
    document.body.insertBefore(message, document.body.firstChild);
    
    setTimeout(() => {
        if (message.parentNode) {
            message.parentNode.removeChild(message);
        }
    }, 3000);
}

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    // Автоматически обновляем список получателей при загрузке
    if (currentUser) {
        updateRecipientsList();
    }
});