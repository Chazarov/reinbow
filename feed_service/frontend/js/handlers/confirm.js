import FeedAPI from '../client/feed_service'
/**
 * Инициализирует обработчик кнопки повторной отправки письма подтверждения (`#repeat-confirm`).
 *
 * При клике отправляет GET-запрос на `/users/me/confirmation` с ID сохранённого пользователя.
 * В случае ошибки выводит сообщение в элемент `#confirm-error`.
 *
 * @function initRepeatConfirmHandlers
 * @requires HTML-элементы с id: `repeat-confirm`, `confirm-error`.
 * @returns {void}
 *
 * @example
 * Вызов после загрузки DOM
 * document.addEventListener('DOMContentLoaded', initRepeatConfirmHandlers);
 */
export function initRepeatConfirmHandlers() {
    const confirmBtn = document.getElementById('repeat-confirm');
    const confirmError = document.getElementById('confirm-error');

    confirmBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        confirmError.textContent = '';
        confirmBtn.disabled = true;

        try {
            await FeedAPI.sendConfirmation();
        } catch (err) {
            confirmError.textContent = err.message;
        } finally {
            confirmBtn.disabled = false;
        }
    });
}