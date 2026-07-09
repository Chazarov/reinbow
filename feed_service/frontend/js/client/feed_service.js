import {showGuestUI} from "../index";

class ApiError extends Error {
    constructor(status, message) {
        super(message);
        this.status = status;
        this.name = 'ApiError';
    }
}


class FeedServiceClient {
    async request(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(response.status, 'Ошибка сети');
        const data = await response.json();
        if (!data.success) throw new Error(400, data.err_message || 'Ошибка сервера');
        return data.data;
    }

    async getUserData() {
        try {
            const data = await this.request("/api/users/me");
            if (data?.user?.data) {
                data.user.data.is_confirmed = data.user.data.is_confirmed === "true";
            }
            return data.user;
        } catch {
            return undefined;
        }
    }

    async getMainFeed() {
        return this.request('/api/feed');
    }
    async getUserFeed(userId) {
        return this.request(`/api/feed?user_id=${userId}`);
    }

    async logout() {
        try {
            await fetch('/api/auth/logout', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            });
        } catch {
            // ошибка запроса игнорируется
        }
    }

    /**
     * Загружает пост на сервер.
     * @param {FormData} formData - данные формы с файлом и описанием
     * @returns {Promise<Object>} ответ сервера (например, созданный пост)
     * @throws {ApiError} с status=401 при неавторизованном доступе
     * @throws {Error} при других ошибках загрузки
     */
    async upload(formData) {
        const response = await fetch('/api/posts', {
            method: 'POST',
            body: formData,
        });
        const data = await response.json();

        if (response.status === 401) {
            throw new ApiError(401, 'Требуется авторизация');
        }

        if (!response.ok || !data.success) {
            throw new Error(data.err_message || 'Ошибка загрузки');
        }

        return data.data;
    }

    /**
     * Отправляет запрос на повторную отправку письма подтверждения.
     * @returns {Promise<Object>} данные ответа сервера
     * @throws {Error} при сетевой ошибке или ошибке сервера
     */
    async sendConfirmation() {
        const response = await fetch('/api/users/me/confirmation', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        });
        const data = await response.json();
        if (!response.ok || !data.success) {
            throw new Error(data.err_message || 'Ошибка. Попробуйте позже');
        }
        return data;
    }

    /**
     * Регистрация нового пользователя.
     * @param {string} username
     * @param {string} email
     * @param {string} password
     * @returns {Promise<{user: Object}>} данные созданного пользователя
     * @throws {Error} при ошибке регистрации или некорректном ответе
     */
    async signup(username, email, password) {
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });
        const data = await response.json();
        if (!response.ok || !data.success) {
            throw new Error(data.err_message || 'Ошибка регистрации');
        }
        const { user } = data.data;
        if (!user) {
            throw new Error('Некорректный ответ сервера');
        }
        return { user };
    }

    /**
     * Аутентификация пользователя.
     * @param {string} username
     * @param {string} password
     * @returns {Promise<{user: Object}>} объект с полем user (и опционально токенами)
     * @throws {Error} при неверных данных или сетевой ошибке
     */
    async signin(username, password) {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
            credentials: 'include'
        });
        const data = await response.json();
        if (!response.ok || !data.success) {
            throw new Error(data.err_message || 'Ошибка входа');
        }
        const { access_token, refresh_token, user } = data.data;
        if (!access_token || !refresh_token || !user) {
            throw new Error('Некорректный ответ сервера');
        }
        return { access_token, refresh_token, user };
    }
}

export default new FeedServiceClient();
export { ApiError };