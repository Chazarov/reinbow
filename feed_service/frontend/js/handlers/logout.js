import { showGuestUI } from '../index.js';
import FeedAPI from '../client/feed_service.js';

/**
 * Инициализирует обработчик клика по кнопке выхода (`#btnLogout`).
 *
 * При клике:
 * - Отправляет POST-запрос на `/api/auth/logout` через клиент API
 *   (ошибка запроса игнорируется).
 * - Обновляет интерфейс вызовом {@link module:main.showGuestUI}.
 *
 * @function initLogoutHandler
 * @requires module:main.showGuestUI
 * @requires FeedAPI
 * @requires HTML-элемент с id: `btnLogout`.
 * @returns {void}
 */
export function initLogoutHandler() {
    const btnLogout = document.getElementById('btnLogout');

    btnLogout.addEventListener('click', async () => {
        await FeedAPI.logout();   // инвалидация сессии, ошибка игнорируется
        showGuestUI();
    });
}