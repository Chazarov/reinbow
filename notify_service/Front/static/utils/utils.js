export const predefinedTypes = ["user_register", "user_login", "admin_newImg", "user_imgVerdict"];

export const emailIdMap = {
    'user_register': 'email_reg',
    'user_login': 'email_login',
    'user_imgVerdict': 'email_user_imgVerdict',
    'admin_newImg': 'email_admin_newImg'
};

export const tgIdMap = {
    'user_register': 'tg_reg',
    'user_login': 'tg_login',
    'user_imgVerdict': 'tg_user_imgVerdict',
    'admin_newImg': 'tg_admin_newImg'
};

export const defaultNotifyNames = {
    'user_register': 'Регистрация аккаунта',
    'user_login': 'Вход в аккаунт',
    'user_imgVerdict': 'Решение по проверке фото',
    'admin_newImg': '[admin] Новое фото для модерации'
};

export let disabledPredefinedTypes = new Set(JSON.parse(localStorage.getItem('disabledPredefinedTypes') || '[]'));

export function saveDisabledTypes() {
    localStorage.setItem('disabledPredefinedTypes', JSON.stringify([...disabledPredefinedTypes]));
}

export function getTableBody() {
    const tbody = document.getElementById('table_main_body');
    if (tbody) return tbody;

    const table = document.querySelector('.table_main');
    if (!table) return null;
    let tbodyEl = table.querySelector('tbody');
    if (!tbodyEl) {
        tbodyEl = document.createElement('tbody');
        tbodyEl.id = 'table_main_body';
        table.appendChild(tbodyEl);
    }
    return tbodyEl;
}