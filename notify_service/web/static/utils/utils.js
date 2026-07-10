/**
 * Массив предопределенных типов уведомлений
 * @type {string[]}
 */
export const predefinedTypes = ["user_register", "user_login", "admin_newImg", "user_imgVerdict"];

/**
 * Объект, сопоставляющий типы уведомлений с ID элементов email чекбоксов
 * @type {Object<string, string>}
 */
export const emailIdMap = {
    'user_register': 'email_reg',
    'user_login': 'email_login',
    'user_imgVerdict': 'email_user_imgVerdict',
    'admin_newImg': 'email_admin_newImg'
};

/**
 * Объект, сопоставляющий типы уведомлений с ID элементов telegram чекбоксов
 * @type {Object<string, string>}
 */
export const tgIdMap = {
    'user_register': 'tg_reg',
    'user_login': 'tg_login',
    'user_imgVerdict': 'tg_user_imgVerdict',
    'admin_newImg': 'tg_admin_newImg'
};

/**
 * Объект, содержащий имена по умолчанию для типов уведомлений
 * @type {Object<string, string>}
 */
export const defaultNotifyNames = {
    'user_register': 'Регистрация аккаунта',
    'user_login': 'Вход в аккаунт',
    'user_imgVerdict': 'Решение по проверке фото',
    'admin_newImg': '[admin] Новое фото для модерации'
};

/**
 * Набор отключенных типов уведомлений, сохраняемый в localStorage
 * @type {Set<string>}
 */
export let disabledPredefinedTypes = new Set(JSON.parse(localStorage.getItem('disabledPredefinedTypes') || '[]'));

/**
 * Сохраняет текущий набор отключенных типов уведомлений в localStorage
 * Преобразует Set в массив и сохраняет его как JSON строку
 */
export function saveDisabledTypes() {
    localStorage.setItem('disabledPredefinedTypes', JSON.stringify([...disabledPredefinedTypes]));
}

/**
 * Получает элемент tbody таблицы настроек
 * Если элемент tbody с id 'table_main_body' существует, возвращает его
 * Если таблица существует, но tbody отсутствует, создает новый tbody элемент
 * @returns {HTMLTableSectionElement|null} - Элемент tbody таблицы или null, если таблица не найдена
 */
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