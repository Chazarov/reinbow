import '../css/style.css';
import { initFeed } from './feed/feed.js';
import { initLogoutHandler } from './handlers/logout.js';
import { initSigninHandlers } from './handlers/signin.js';
import { initSignupHandlers } from './handlers/signup.js';
import { initUploadHandlers } from './handlers/upload.js';

document.addEventListener('DOMContentLoaded', () => {
    // Модальные окна
    const signinModal = document.getElementById('signinModal');
    const signupModal = document.getElementById('signupModal');
    const btnSignin = document.getElementById('btnSignin');
    const btnSignup = document.getElementById('btnSignup');
    const switchToSignup = document.getElementById('switchToSignup');
    const switchToSignin = document.getElementById('switchToSignin');
    const closeButtons = document.querySelectorAll('.close');

    // DOM-элементы для управления UI
    const guestButtons = document.getElementById('guestButtons');
    const userBlock = document.getElementById('userBlock');
    const userNameDisplay = document.getElementById('userNameDisplay');

    // Сброс формы и сообщения об ошибке
    window.resetForm = function(form, errorElement) {
        form.reset();
        if (errorElement) errorElement.textContent = '';
    };

    // Открытие модального окна
    window.openModal = function(modal) {
        modal.classList.add('active');
    };

    // Закрытие модального окна со сбросом формы
    window.closeModal = function(modal) {
        modal.classList.remove('active');
        // сброс всех форм в модалке
        const forms = modal.querySelectorAll('form');
        forms.forEach(f => f.reset());
        // очистка ошибок
        modal.querySelectorAll('.error-message').forEach(el => el.textContent = '');
        // очистка списка файлов (если есть)
        const fileList = modal.querySelector('#fileList');
        if (fileList) fileList.innerHTML = '';
    };

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

    // Закрытие по клику вне окна
    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            closeModal(e.target);
        }
    });

    // Получаем сохранённого пользователя в виде Object
    window.getSavedUser = () => {
        const userStr = localStorage.getItem('user');
        return userStr ? JSON.parse(userStr) : null;
    }

    // Обновляем интерфейс: показываем блок пользователя, скрываем гостевые кнопки
    window.showLoggedInUI = (user) => {
        if (guestButtons) guestButtons.style.display = 'none';
        if (userBlock) userBlock.style.display = 'block';
        if (userNameDisplay) userNameDisplay.textContent = user.username;
    }

    // Обновляем интерфейс для гостя
    window.showGuestUI = () => {
        if (guestButtons) guestButtons.style.display = 'block';
        if (userBlock) userBlock.style.display = 'none';
        if (userNameDisplay) userNameDisplay.textContent = '';
    }

    // При загрузке страницы проверяем, есть ли сохранённая сессия
    function checkAuthOnLoad() {
        const user = getSavedUser();
        if (user) {
            showLoggedInUI(user);
        } else {
            showGuestUI();
        }
    }

    // Вызов при старте
    initSignupHandlers();

    initSigninHandlers();

    initLogoutHandler();

    checkAuthOnLoad();

    initFeed();

    initUploadHandlers();
});