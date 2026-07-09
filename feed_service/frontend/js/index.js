/**
 * Главный модуль приложения (точка входа).
 *
 * Импортирует стили и все необходимые модули. При загрузке DOM инициализирует:
 * - модальные окна входа и регистрации,
 * - экспортируемые функции управления UI (`openModal`, `closeModal`, `resetForm`, `showLoggedInUI`, `showGuestUI`, `toggleConfirmedUI`),
 * - обработчики форм (вход, регистрация, выход, загрузка постов),
 * - проверку текущей сессии и обновление интерфейса,
 * - загрузку основной ленты.
 *
 * Экспортируемые функции:
 * - {@link openModal} — открыть модальное окно.
 * - {@link closeModal} — закрыть модальное окно со сбросом форм и очисткой ошибок.
 * - {@link resetForm} — сбросить форму и очистить элемент ошибки.
 * - {@link showLoggedInUI} — переключить интерфейс на авторизованного пользователя.
 * - {@link showGuestUI} — переключить интерфейс на гостя.
 * - {@link toggleConfirmedUI} — обновить UI в зависимости от статуса подтверждения учётной записи.
 *
 * @file index.js
 * @module main
 * @requires module:feed/feed
 * @requires module:handlers/logout
 * @requires module:handlers/signin
 * @requires module:handlers/signup
 * @requires module:handlers/upload
 * @requires module:handlers/confirm
 */

import '../css/style.css';
import '../css/feed.css';
import '../css/form.css';
import '../css/modal.css';
import '../css/font-awesome.min.css';
import {initFeed, loadMainFeed} from './feed/feed.js';
import { initLogoutHandler } from './handlers/logout.js';
import { initSigninHandlers } from './handlers/signin.js';
import { initSignupHandlers } from './handlers/signup.js';
import { initUploadHandlers } from './handlers/upload.js';
import { initRepeatConfirmHandlers } from './handlers/confirm.js';
import FeedAPI from "./client/feed_service.js";

document.addEventListener('DOMContentLoaded', () => {
    // Модальные окна
    const signinModal = document.getElementById('signinModal');
    const signupModal = document.getElementById('signupModal');
    const btnSignin = document.getElementById('btnSignin');
    const btnSignup = document.getElementById('btnSignup');
    const switchToSignup = document.getElementById('switchToSignup');
    const switchToSignin = document.getElementById('switchToSignin');
    const closeButtons = document.querySelectorAll('.close');

    // Обработчики открытия
    btnSignin.addEventListener('click', () => openModal(signinModal));
    btnSignup.addEventListener('click', () => openModal(signupModal));

    // Переключение между окнами
    switchToSignup.addEventListener('click', (e) => {
        e.preventDefault();
        closeModal(signinModal);
        openModal(signupModal);
    });

    switchToSignin.addEventListener('click', (e) => {
        e.preventDefault();
        closeModal(signupModal);
        openModal(signinModal);
    });

    // Закрытие по крестику
    closeButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            const modal = e.target.closest('.modal');
            closeModal(modal);
        });
    });

    let mousedownTarget = null;

    window.addEventListener('mousedown', (e) => {
        mousedownTarget = e.target;
    });

    // Закрытие по клику вне окна
    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal') && mousedownTarget === e.target) {
            closeModal(e.target);
        }
    });

    /**
     * Проверяет наличие сессии при загрузке страницы.
     * Пытается обновить токен доступа, затем на основе наличия сохранённого
     * пользователя показывает соответствующий интерфейс.
     *
     * @inner
     * @function checkAuthOnLoad
     */
    async function checkAuthOnLoad() {
        const user = await FeedAPI.getUserData();
        if (!user) {
            showGuestUI();
            return;
        }
        showLoggedInUI(user);
        toggleConfirmedUI();
    }

    // Вызов при старте
    initSignupHandlers();
    initSigninHandlers();
    initLogoutHandler();
    checkAuthOnLoad();
    initFeed();
    initUploadHandlers();
    initRepeatConfirmHandlers();

    const hash = window.location.hash;
    const value = hash.slice(1);
    if(value === "login") {
        openModal(signinModal)
    }
});

/**
 * Сбрасывает форму и очищает элемент с сообщением об ошибке.
 *
 * @function resetForm
 * @param {HTMLFormElement} form - форма для сброса.
 * @param {HTMLElement} [errorElement] - элемент для очистки текста ошибки.
 * @returns {void}
 */
export function resetForm(form, errorElement) {
    form.reset();
    if (errorElement) errorElement.textContent = '';
}

/**
 * Открывает модальное окно, добавляя класс `active`.
 *
 * @function openModal
 * @param {HTMLElement} modal - DOM-элемент модального окна.
 * @returns {void}
 */
export function openModal(modal) {
    modal.classList.add('active');
}

/**
 * Закрывает модальное окно, удаляя класс `active`.
 * Дополнительно сбрасывает все формы внутри модалки, очищает все
 * элементы с классом `error-message` и список файлов `#fileList`.
 *
 * @function closeModal
 * @param {HTMLElement} modal - DOM-элемент модального окна.
 * @returns {void}
 */
export function closeModal(modal) {
    modal.classList.remove('active');
    // сброс всех форм в модалке
    const forms = modal.querySelectorAll('form');
    forms.forEach(f => f.reset());
    // очистка ошибок
    modal.querySelectorAll('.error-message').forEach(el => el.textContent = '');
    // очистка списка файлов (если есть)
    const fileList = modal.querySelector('#fileList');
    if (fileList) fileList.innerHTML = '';
}

/**
 * Переключает интерфейс на отображение для авторизованного пользователя.
 * Скрывает гостевые кнопки, показывает блок пользователя с именем и
 * кнопку загрузки.
 *
 * @function showLoggedInUI
 * @param {Object} user - объект пользователя с полем `username`.
 * @param {string} user.username - отображаемое имя.
 * @returns {void}
 */
export function showLoggedInUI(user) {
    const guestButtons = document.getElementById('guestButtons');
    const userBlock = document.getElementById('userBlock');
    const userNameDisplay = document.getElementById('user-name-display');
    const uploadBtn = document.getElementById('btn-upload');
    const confirmBtn = document.getElementById('repeat-confirm');

    if (guestButtons) guestButtons.style.display = 'none';
    if (userBlock) userBlock.style.display = 'block';
    if (userNameDisplay) userNameDisplay.textContent = user.login;
    if (uploadBtn) uploadBtn.style.display = 'inline-flex';
    if (confirmBtn) uploadBtn.style.display = 'block';
}

/**
 * Переключает интерфейс в гостевой режим.
 * Показывает кнопки входа/регистрации, скрывает блок пользователя и
 * кнопку загрузки.
 *
 * @function showGuestUI
 * @returns {void}
 */
export function showGuestUI() {
    const guestButtons = document.getElementById('guestButtons');
    const userBlock = document.getElementById('userBlock');
    const userNameDisplay = document.getElementById('user-name-display');
    const uploadBtn = document.getElementById('btn-upload');
    const confirmBtn = document.getElementById('repeat-confirm');

    if (guestButtons) guestButtons.style.display = 'block';
    if (userBlock) userBlock.style.display = 'none';
    if (userNameDisplay) userNameDisplay.textContent = '';
    if (uploadBtn) uploadBtn.style.display = 'none';
    if (confirmBtn) confirmBtn.style.display = 'none';
}

/**
 * Показывает элементы UI, доступные только при неактивной учетной записи (is_confirmed = false).
 * Подсвечивает никнейм красным цветом и управляет видимостью кнопки добавления постов.
 *
 * @async
 * @function toggleConfirmedUI
 * @returns {Promise<void>}
 */
export async function toggleConfirmedUI() {
    const user = await FeedAPI.getUserData();

    if (!user) {
        showGuestUI();
        loadMainFeed();
        return;
    }

    const userNameDisplay = document.getElementById('user-name-display');
    const uploadBtn = document.getElementById('btn-upload');

    if (userNameDisplay && uploadBtn) {
        user.data.is_confirmed
            ? userNameDisplay.classList.remove("not-confirmed")
            : userNameDisplay.classList.add("not-confirmed");
        user.data.is_confirmed
            ? uploadBtn.style.display = 'inline-flex'
            : uploadBtn.style.display = 'none';
    }
}