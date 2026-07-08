import './auth.css';
import * as api from '../api.js';
import { showNotification } from '../notifications/notifications.js';
import { GetSettings, hideDisabledRows } from '../settings/settings.js';
import { initAllMultiselects, webhookUrls, webhookSelections, multiselectInstances, renderWebhookTable } from '../webhooks/webhooks.js';

export let isLogged = false;

export function closeModal() {
    let modal_window = document.getElementById('modalWindow');
    let black_background = document.getElementById('blackBackground');
    modal_window.classList.remove('open');
    modal_window.classList.add('closed');
    black_background.classList.remove('open');
    black_background.classList.add('closed');
}

export function openModal() {
    let modal_window = document.getElementById('modalWindow');
    let black_background = document.getElementById('blackBackground');
    modal_window.classList.remove('closed');
    modal_window.classList.add('open');
    black_background.classList.remove('closed');
    black_background.classList.add('open');
}

export async function loginProcedure() {
    let login_value = document.getElementById('login_field').value.trim();
    let password_value = document.getElementById('password_field').value.trim();

    if (!login_value || !password_value) {
        showNotification('error', 'Ошибка', 'Заполните все поля');
        return;
    }

    const payload = {
        "login": login_value,
        "password": password_value
    }

    try {
        const result = await api.login(payload);
        console.log(result);
        if (result.success == true) {
            console.log("Вход успешен!");
            localStorage.setItem('token', result.token);
            isLogged = true;
            document.getElementById('not_Logged').classList.add('auth-hidden');
            document.getElementById('isAuthorized').classList.remove('auth-hidden');
            closeModal();
            showNotification('success', 'Успех!', 'Успешный вход');
            GetSettings().then(() => hideDisabledRows());
        } else {
            console.log("Неправильные данные!");
            showNotification('error', 'Ошибка', result.Error_message || 'Неверный логин или пароль');
        }
    } catch (error) {
        console.error('Login error:', error);
        showNotification('error', 'Ошибка', 'Ошибка при попытке входа');
    }
}

export async function logoutProcedure() {
    localStorage.removeItem('token');
    isLogged = false;

    document.getElementById('login_field').value = "";
    document.getElementById('password_field').value = "";

    document.getElementById('not_Logged').classList.remove('auth-hidden');
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

export function restoreAuthState() {
    const token = localStorage.getItem('token');

    if (token) {
        isLogged = true;
        document.getElementById('not_Logged').classList.add('auth-hidden');
        document.getElementById('isAuthorized').classList.remove('auth-hidden');
        console.log("Добро пожаловать, admin");
        initAllMultiselects();
        GetSettings().then(() => hideDisabledRows());
    } else {
        isLogged = false;
        document.getElementById('not_Logged').classList.remove('auth-hidden');
        document.getElementById('isAuthorized').classList.add('auth-hidden');
    }
}

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