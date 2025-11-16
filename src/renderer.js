let currentUser = null;
let authToken = null;

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    checkAuthStatus();
    updateRecipientsList();
});

// Проверка статуса авторизации
function checkAuthStatus() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (token && user) {
        authToken = token;
        currentUser = JSON.parse(user);
        showMainSection();
        loadUserData();
    }
}

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
        const result = await window.electronAPI.register({
            cipher, password, name, group
        });
        
        if (result.token && result.user) {
            // Сохраняем токен и данные пользователя
            localStorage.setItem('token', result.token);
            localStorage.setItem('user', JSON.stringify(result.user));
            authToken = result.token;
            currentUser = result.user;
            
            showMessage(result.message, 'success');
            showMainSection();
            loadUserData();
            
            // Очищаем форму
            document.getElementById('reg-cipher').value = '';
            document.getElementById('reg-password').value = '';
            document.getElementById('reg-name').value = '';
            document.getElementById('reg-group').value = '';
        }
    } catch (error) {
        showMessage(error.message, 'error');
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
        const result = await window.electronAPI.login(cipher, password);
        
        if (result.token && result.user) {
            // Сохраняем токен и данные пользователя
            localStorage.setItem('token', result.token);
            localStorage.setItem('user', JSON.stringify(result.user));
            authToken = result.token;
            currentUser = result.user;
            
            showMainSection();
            loadUserData();
        }
    } catch (error) {
        showMessage(error.message, 'error');
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
        const [profileData, transactionsData] = await Promise.all([
            window.electronAPI.getUserProfile(),
            window.electronAPI.getUserTransactions()
        ]);
        
        // Обновляем данные пользователя
        currentUser = { ...currentUser, ...profileData.user };
        localStorage.setItem('user', JSON.stringify(currentUser));
        
        document.getElementById('balance-amount').textContent = 
            `${currentUser.balance} Кибиков`;
        document.getElementById('user-name').textContent = currentUser.name;
        
        updateTransactionsList(transactionsData.transactions);
        loadTasks();
    } catch (error) {
        showMessage('Ошибка загрузки данных', 'error');
        if (error.message.includes('авторизация')) {
            logout();
        }
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
                <div class="transaction-name">${transaction.recipientName}</div>
                <div class="transaction-date">${new Date(transaction.date).toLocaleDateString('ru-RU')}</div>
            </div>
            <div class="transaction-amount ${transaction.toCipher === currentUser.cipher ? 'positive' : 'negative'}">
                ${transaction.toCipher === currentUser.cipher ? '+' : '-'}${transaction.amount} К
            </div>
        </div>
    `).join('');
}

// Обновление списка получателей
async function updateRecipientsList() {
    const select = document.getElementById('transfer-recipient');
    select.innerHTML = '<option value="">Выберите получателя</option>';
    
    // Тестовые пользователи для демонстрации
    const testUsers = [
        { cipher: 'kkso07_001', name: 'Иван Петров (ККСО-07-23)', group: 'кксо-07-23' },
        { cipher: 'kkso07_002', name: 'Мария Сидорова (ККСО-07-23)', group: 'кксо-07-23' },
        { cipher: 'kkso06_001', name: 'Алексей Козлов (ККСО-06-23)', group: 'кксо-06-23' },
        { cipher: 'kkso06_002', name: 'Елена Новикова (ККСО-06-23)', group: 'кксо-06-23' }
    ];
    
    testUsers.forEach(user => {
        const option = document.createElement('option');
        option.value = user.cipher;
        option.textContent = user.name;
        option.dataset.group = user.group;
        select.appendChild(option);
    });
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
    const amount = parseInt(document.getElementById('transfer-amount').value);
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
        const result = await window.electronAPI.makeTransfer({
            toCipher: toCipher,
            amount: amount,
            recipientName: recipientOption.textContent.split(' (')[0]
        });
        
        showMessage(result.message, 'success');
        document.getElementById('transfer-amount').value = '';
        loadUserData(); // Обновляем данные
    } catch (error) {
        showMessage(error.message, 'error');
    }
}

// Загрузка заданий
async function loadTasks() {
    try {
        const result = await window.electronAPI.getTasks();
        const container = document.getElementById('tasks-container');
        
        // Временная реализация - задания всегда не выполнены
        // В реальном приложении нужно добавить проверку статуса выполнения
        container.innerHTML = result.tasks.map(task => `
            <div class="task-item">
                <div class="task-info">
                    <div class="task-title">${task.title}</div>
                    <div class="task-status">⏳ Доступно</div>
                </div>
                <div class="task-reward">+${task.reward} К</div>
                <button onclick="completeTask(${task.id})" class="btn-primary" style="margin-left: 10px; padding: 8px 16px;">
                    Выполнить
                </button>
            </div>
        `).join('');
    } catch (error) {
        showMessage('Ошибка загрузки заданий', 'error');
    }
}

// Выполнение задания
async function completeTask(taskId) {
    if (!currentUser) return;

    try {
        const result = await window.electronAPI.completeTask(taskId);
        showMessage(result.message, 'success');
        loadUserData(); // Обновляем баланс и задания
    } catch (error) {
        showMessage(error.message, 'error');
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
    authToken = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
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