/**
 * Инициализирует обработчик отправки формы входа (`#signinForm`).
 *
 * Выполняет клиентскую валидацию: поля логина и пароля обязательны.
 *
 * При успешном входе:
 * - обновляет интерфейс вызовом {@link module:main.showLoggedInUI},
 * - закрывает модальное окно через {@link module:main.closeModal}
 *   (предполагается, что форма находится в модальном окне `#signinModal`).
 *
 * Ошибки (сетевые, API, отсутствие обязательных данных в ответе)
 * выводятся в элемент `#signinError`.
 *
 * @function initSigninHandlers
 * @requires module:main.closeModal
 * @requires module:main.showLoggedInUI
 * @requires module:main.toggleConfirmedUI
 * @requires HTML-элементы с id: `signinForm`, `signinError`, `signinUsername`,
 *           `signinPassword`, `signinModal`.
 * @returns {void}
 *
 * @example
 * Вызов после загрузки DOM
 * document.addEventListener('DOMContentLoaded', initSigninHandlers);
 */
import { closeModal, showLoggedInUI, toggleConfirmedUI } from '../index.js';
import FeedAPI from '../client/feed_service'

export function initSigninHandlers() {
    const signinForm = document.getElementById('signinForm');
    const signinError = document.getElementById('signinError');
    const signinModal = document.getElementById('signinModal');

    signinForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        signinError.textContent = '';
        e.submitter.disabled = true;

        const username = document.getElementById('signinUsername').value.trim();
        const password = document.getElementById('signinPassword').value;

        if (!username || !password) {
            signinError.textContent = 'Заполните все поля';
            e.submitter.disabled = false;
            return;
        }

        try {
            const { user } = await FeedAPI.signin(username, password);
            showLoggedInUI(user);
            closeModal(signinModal);
            toggleConfirmedUI();
        } catch (err) {
            signinError.textContent = err.message;
        } finally {
            e.submitter.disabled = false;
        }
    });
}