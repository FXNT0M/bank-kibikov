const { contextBridge } = require('electron');

// API для работы с бэкендом
const api = {
  // Базовый URL бэкенда
  baseURL: 'http://localhost:8080/api/v1',
  
  // Общая функция для HTTP запросов
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    if (config.body && typeof config.body === 'object') {
      config.body = JSON.stringify(config.body);
    }

    try {
      const response = await fetch(url, config);
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Ошибка сервера');
      }
      
      return data;
    } catch (error) {
      throw error;
    }
  },

  // Авторизация
  async register(userData) {
    return this.request('/auth/register', {
      method: 'POST',
      body: userData,
    });
  },

  async login(cipher, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: { cipher, password },
    });
  },

  // Защищенные запросы
  async authenticatedRequest(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('Требуется авторизация');
    }

    return this.request(endpoint, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        ...options.headers,
      },
    });
  },

  // API методы
  async getUserProfile() {
    return this.authenticatedRequest('/user/profile', {
      method: 'GET',
    });
  },

  async getUserTransactions() {
    return this.authenticatedRequest('/user/transactions', {
      method: 'GET',
    });
  },

  async makeTransfer(transferData) {
    return this.authenticatedRequest('/transactions/transfer', {
      method: 'POST',
      body: transferData,
    });
  },

  async getTasks() {
    return this.authenticatedRequest('/tasks', {
      method: 'GET',
    });
  },

  async completeTask(taskId) {
    return this.authenticatedRequest(`/tasks/${taskId}/complete`, {
      method: 'POST',
    });
  }
};

// Предоставляем API в рендерер
contextBridge.exposeInMainWorld('electronAPI', api);