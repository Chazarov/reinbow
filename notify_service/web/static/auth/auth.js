import './auth.css';
import * as api from '../api.js';
import { showNotification } from '../notifications/notifications.js';
import { GetSettings, hideDisabledRows } from '../settings/settings.js';
import { initAllMultiselects, webhookUrls, webhookSelections, multiselectInstances, renderWebhookTable } from '../webhooks/webhooks.js';

export let isLogged = false;

/**
 * Закрывает модальное окно аутентификации
 * Добавляет класс 'closed' и удаляет класс 'open' у модального окна и затемненного фона
 */
export function closeModal() {
    let modalWindow = document.getElementById('modalWindow');
    let blackBackground = document.getElementById('blackBackground');
    modalWindow.classList.remove('open');
    modalWindow.classList.add('closed');
    blackBackground.classList.remove('open');
    blackBackground.classList.add('closed');
}

/**
 * Открывает модальное окно аутентификации
 * Добавляет класс 'open' и удаляет класс 'closed' у модального окна и затемненного фона
 */
export function openModal() {
    let modalWindow = document.getElementById('modalWindow');
    let blackBackground = document.getElementById('blackBackground');
    modalWindow.classList.remove('closed');
    modalWindow.classList.add('open');
    blackBackground.classList.remove('closed');
    blackBackground.classList.add('open');
}

/**
 * Обрабатывает процесс входа пользователя в систему
 * Получает значения логина и пароля из полей ввода, отправляет запрос на аутентификацию
 * При успешном входе сохраняет токен в localStorage, обновляет UI и показывает уведомление об успехе
 * При неудаче показывает сообщение об ошибке
 */
export async function loginProcedure() {
    let loginValue = document.getElementById('loginField').value.trim();
    let passwordValue = document.getElementById('passwordField').value.trim();

    if (!loginValue || !passwordValue) {
        showNotification('error', 'Ошибка', 'Заполните все поля');
        return;
    }

    const payload = {
        "login": loginValue,
        "pass": passwordValue
    }

    try {
        const result = await api.login(payload);
        console.log(result);
        if (result.success == true) {
            console.log("Вход успешен!");
            localStorage.setItem('token', result.token);
            isLogged = true;
            document.getElementById('notLogged').classList.add('auth-hidden');
            document.getElementById('isAuthorized').classList.remove('auth-hidden');
            closeModal();
            showNotification('success', 'Успех!', 'Успешный вход');
            GetSettings().then(() => hideDisabledRows());
        } else {
            console.log("Неправильные данные!");
            showNotification('error', 'Ошибка', result.ErrorMessage || 'Неверный логин или пароль');
        }
    } catch (error) {
        console.error('Login error:', error);
        showNotification('error', 'Ошибка', 'Ошибка при попытке входа');
    }
}

/**
 * Обрабатывает процесс выхода пользователя из системы
 * Очищает токен из localStorage, сбрасывает состояние авторизации,
 * очищает форму входа, обновляет UI и сбрасывает настройки вебхуков
 */
export async function logoutProcedure() {
    try {
        await api.logout(); // сервер удалит cookie
    } catch (error) {
        console.error('Logout error:', error);
    }

    //localStorage.removeItem('token');
    isLogged = false;

    document.getElementById('loginField').value = "";
    document.getElementById('passwordField').value = "";

    document.getElementById('notLogged').classList.remove('auth-hidden');
    document.getElementById('isAuthorized').classList.add('auth-hidden');

    webhookUrls.length = 0;
    Object.keys(webhookSelections).forEach(key => delete webhookSelections[key]);
    multiselectInstances.length = 0;
    initAllMultiselects();
    renderWebhookTable();

    console.log("Выход успешен!");
    closeModal();
    showNotification('success', 'Успех', 'Выход из аккаунта успешен!');
}

export async function restoreAuthState() {
    try {
        const userData = await api.getMe();
        // Если запрос успешен, считаем пользователя авторизованным
        isLogged = true;
        currentUser = userData; // сохраняем данные
        updateUIForLoggedIn();
        console.log('Сессия восстановлена, пользователь:', userData);
        // Дополнительно инициализируем остальные части приложения
        initAllMultiselects();
        await GetSettings();
        hideDisabledRows();
    } catch (error) {
        // Ошибка (401 или другая) – пользователь не авторизован
        console.warn('Не удалось восстановить сессию:', error);
        isLogged = false;
        currentUser = null;
        updateUIForLoggedOut();
    }
}

// Вспомогательные функции для обновления UI
function updateUIForLoggedIn() {
    const notLogged = document.getElementById('notLogged');
    const isAuthorized = document.getElementById('isAuthorized');
    if (notLogged) notLogged.classList.add('auth-hidden');
    if (isAuthorized) isAuthorized.classList.remove('auth-hidden');

    // Обновляем имя пользователя, если есть
    if (currentUser && currentUser.login) {
        const userNameElement = document.querySelector('.userAuthorized p');
        if (userNameElement) userNameElement.textContent = currentUser.login;
    }

    // Инициализация остальных частей приложения
    initAllMultiselects();
    GetSettings().then(() => hideDisabledRows());
}

/**
 * Восстанавливает состояние аутентификации при загрузке страницы
 * Проверяет наличие токена в localStorage и восстанавливает состояние авторизации
 */
/**
 * Восстанавливает состояние аутентификации при загрузке страницы
 * Проверяет наличие токена в localStorage и восстанавливает состояние авторизации
 */
// export function restoreAuthState() {
//     const token = localStorage.getItem('token');

//     if (token) {
//         isLogged = true;
//         document.getElementById('notLogged').classList.add('auth-hidden');
//         document.getElementById('isAuthorized').classList.remove('auth-hidden');
//         console.log("Добро пожаловать, admin");
//         initAllMultiselects();
//         GetSettings().then(() => hideDisabledRows());
//     } else {
//         isLogged = false;
//         document.getElementById('notLogged').classList.remove('auth-hidden');
//         document.getElementById('isAuthorized').classList.add('auth-hidden');
//     }
// }

/**
 * Инициализирует обработчики событий для элементов аутентификации
 * Назначает обработчики кликов на кнопки открытия/закрытия модального окна,
 * входа и выхода из системы
 */
export function initAuthListeners() {
    const openModalBtn = document.getElementById('openModalBtn');
    if (openModalBtn) openModalBtn.addEventListener('click', openModal);

    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) logoutBtn.addEventListener('click', logoutProcedure);

    const closeModalBtn = document.getElementById('closeModalBtn');
    if (closeModalBtn) closeModalBtn.addEventListener('click', closeModal);

    const loginBtn = document.getElementById('loginBtn');
    if (loginBtn) loginBtn.addEventListener('click', loginProcedure);
}