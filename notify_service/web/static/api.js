/**
 * Добавляет токен авторизации к заголовкам запроса
 * @param {string} url - URL для запроса
 * @param {Object} options - Параметры запроса fetch
 * @returns {Promise} - Promise от fetch запроса
 */
export function fetchWithAuth(url, options = {}) {
    // const token = localStorage.getItem('token');
    // if (token) {
    //     options.headers = {
    //         ...options.headers,
    //         'Authorization': `Bearer ${token}`
    //     };
    // }
    return fetch(url, options);
}

/**
 * Выполняет вход пользователя в систему
 * @param {Object} loginData - Объект с данными для входа (логин и пароль)
 * @returns {Promise<Object>} - Promise с результатом входа
 */
export async function login(loginData) {
    const response = await fetch("/api/auth/login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(loginData),
        credentials: 'include'
    });
    return response.json();
}

export async function logout() {
    const response = await fetch("/api/auth/logout", {
        method: 'POST',
        credentials: 'include'
    })
    return response.json();
}
// TODO: add include creditions



// export async function refresh() {
//     const response = await fetch("api/auth/refresh", {
//         method: 'POST'
//     })
//     return response.json();
// }

/**
 * Получает настройки уведомлений с сервера
 * @returns {Promise<Object>} - Promise с данными о настройках уведомлений
 */
export async function getNotifySettings() {
    const response = await fetchWithAuth("/api/notify_settings", {
        credentials: 'include'
    });
    return response.json();
}


export async function getMe() {
    const response = await fetch("/api/auth/me", {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include'
    });
    if (!response.ok) {
        throw new Error(`Auth check failed: ${response.status}`);
    }
    return response.json();
}

// static/api.js

/**
 * Обновляет пару access/refresh токенов через httpOnly куки
 */
export async function refreshTokens() {
    const response = await fetch("/api/auth/refresh", {
        method: 'POST',
        credentials: 'include'   // обязательно!
    });
    if (!response.ok) {
        throw new Error(`Refresh failed: ${response.status}`);
    }
    return response.json();
}

/**
 * Сохраняет настройки уведомлений на сервере
 * @param {Object} payload - Данные настроек для отправки на сервер
 * @returns {Promise<Response>} - Promise с ответом от сервера
 */
export async function saveNotifySettings(payload) {
    return fetchWithAuth("/api/notify_types", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include'
    });
}