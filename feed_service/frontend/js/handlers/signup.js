/**
 * Инициализирует обработчик отправки формы регистрации (`#signupForm`).
 *
 * Выполняет клиентскую валидацию:
 * - Все поля (логин, email, пароль, подтверждение) обязательны.
 * - Пароль и подтверждение должны совпадать.
 * - Email должен соответствовать паттерну `/^[^\s@]+@[^\s@]+\.[^\s@]+$/`.
 *
 * При успешной регистрации:
 * - Закрывает модальное окно регистрации (элемент с id `signupModal`).
 * - Открывает модальное окно подтверждения (элемент с id `confirmModal`).
 * - Обновляет интерфейс вызовом {@link module:main.showLoggedInUI}
 *   и {@link module:main.toggleConfirmedUI}.
 *
 * Ошибки (сетевые, API или валидации) выводятся в элемент `#signupError`.
 *
 * @function initSignupHandlers
 * @requires module:main.closeModal
 * @requires module:main.openModal
 * @requires module:main.showLoggedInUI
 * @requires module:main.toggleConfirmedUI
 * @requires HTML-элементы с id: `signupForm`, `signupError`, `signupUsername`,
 *           `signupEmail`, `signupPassword`, `signupConfirm`,
 *           `signupModal`, `confirmModal`.
 * @returns {void}
 *
 * @example
 * Вызов после загрузки DOM
 * document.addEventListener('DOMContentLoaded', initSignupHandlers);
 */
import { closeModal, openModal, showLoggedInUI, toggleConfirmedUI } from '../index.js';
import FeedAPI from '../client/feed_service'

export function initSignupHandlers() {
    const signupForm = document.getElementById('signupForm');
    const signupError = document.getElementById('signupError');
    const signupModal = document.getElementById('signupModal');
    const confirmModal = document.getElementById('confirmModal');

    signupForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        signupError.textContent = '';
        e.submitter.disabled = true;

        const username = document.getElementById('signupUsername').value.trim();
        const email = document.getElementById('signupEmail').value.trim();
        const password = document.getElementById('signupPassword').value;
        const confirm = document.getElementById('signupConfirm').value;

        if (!username || !email || !password || !confirm) {
            signupError.textContent = 'Заполните все обязательные поля';
            e.submitter.disabled = false;
            return;
        }

        if (password !== confirm) {
            signupError.textContent = 'Пароли не совпадают';
            e.submitter.disabled = false;
            return;
        }

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            signupError.textContent = 'Введите корректный email';
            e.submitter.disabled = false;
            return;
        }

        try {
            const { user } = await FeedAPI.signup(username, email, password);
            // При успехе переключаем модалки и обновляем UI
            closeModal(signupModal);
            openModal(confirmModal);
            showLoggedInUI(user);
            toggleConfirmedUI();
        } catch (err) {
            signupError.textContent = err.message;
        } finally {
            e.submitter.disabled = false;
        }
    });
}