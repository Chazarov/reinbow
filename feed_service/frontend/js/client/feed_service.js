/**
 * Клиентский сервис для взаимодействия с REST API социальной сети.
 * Предоставляет методы для работы с пользователями, лентой постов,
 * аутентификацией и регистрацией.
 *
 * @module FeedServiceClient
 */
import {showGuestUI} from "../index";

/**
 * @typedef {Object} UserProfileData
 * @property {string} email - адрес электронной почты
 * @property {string} created_at - дата регистрации в формате ISO 8601
 * @property {string} is_confirmed - флаг подтверждения почты (строка "true"/"false")
 */

/**
 * @typedef {Object} User
 * @property {number} id - идентификатор пользователя
 * @property {string} login - логин (username)
 * @property {UserProfileData} data - дополнительные данные профиля
 */

/**
 * Ошибка API с HTTP-статусом.
 * Используется для специфичной обработки ошибок авторизации и других
 * ситуаций, когда важен код ответа сервера.
 *
 * @extends Error
 */
class ApiError extends Error {
    /**
     * Создаёт экземпляр ошибки API.
     * @param {number} status - HTTP-статус ответа
     * @param {string} message - Сообщение об ошибке
     */
    constructor(status, message) {
        super(message);
        this.status = status;
        this.name = 'ApiError';
    }
}

/**
 * Основной клиент для запросов к серверу.
 * Все публичные методы возвращают Promise и могут выбрасывать
 * {@link ApiError} или обычную {@link Error}.
 *
 * @class
 */
class FeedServiceClient {
    /**
     * Выполняет GET-запрос к указанному URL и обрабатывает базовые ошибки.
     * Ожидает JSON-ответ с полем `success` и данными в поле `data`.
     *
     * @param {string} url - адрес эндпоинта
     * @returns {Promise<Object>} содержимое поля `data` ответа
     * @throws {Error} при сетевой ошибке или если ответ сервера не успешен
     * @private
     */
    async request(url) {
        const response = await fetch(url);
        if (!response.ok) throw new Error(response.status, 'Ошибка сети');
        const data = await response.json();
        if (!data.success) throw new Error(400, data.err_message || 'Ошибка сервера');
        return data.data;
    }

    /**
     * Получает данные текущего пользователя.
     * Приводит поле `is_confirmed` к булеву типу.
     *
     * @returns {Promise<User|undefined>} объект пользователя или `undefined` при ошибке
     */
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

    /**
     * Загружает основную ленту постов (без параметров).
     *
     * @returns {Promise<Object>} данные ленты
     * @throws {Error} при ошибке сети или сервера
     */
    async getMainFeed() {
        return this.request('/api/feed');
    }

    /**
     * Загружает ленту постов конкретного пользователя.
     *
     * @param {string|number} userId - идентификатор пользователя
     * @returns {Promise<Object>} данные ленты
     * @throws {Error} при ошибке сети или сервера
     */
    async getUserFeed(userId) {
        return this.request(`/api/feed?user_id=${userId}`);
    }

    /**
     * Выполняет выход пользователя (завершение сессии).
     * Ошибки запроса игнорируются, метод всегда завершается без исключений.
     *
     * @returns {Promise<void>}
     */
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
     *
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
     *
     * @returns {Promise<Object>} полный JSON-ответ сервера (содержит `success`, `data` и т.д.)
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
     *
     * @param {string} username - желаемый логин
     * @param {string} email - адрес электронной почты
     * @param {string} password - пароль
     * @returns {Promise<{user: User}>} объект с созданным пользователем
     * @throws {Error} при ошибке регистрации или некорректном ответе
     */
    async signup(username, email, password) {
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login: username, data: { email }, pass: password })
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
     *
     * @param {string} username - логин
     * @param {string} password - пароль
     * @returns {Promise<{user: User}>} объект с данными вошедшего пользователя
     * @throws {Error} при неверных учётных данных или сетевой ошибке
     */
    async signin(username, password) {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login: username, pass: password }),
            credentials: 'include'
        });
        const data = await response.json();
        if (!response.ok || !data.success) {
            throw new Error(data.err_message || 'Ошибка входа');
        }
        const { user } = data.data;
        if (!user) {
            throw new Error('Некорректный ответ сервера');
        }
        return { user };
    }
}

/** Готовый экземпляр клиента для использования в приложении */
export default new FeedServiceClient();

/** Класс ошибки API, позволяющий различать коды ответов */
export { ApiError };