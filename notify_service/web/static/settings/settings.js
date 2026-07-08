import './settings.css';
import * as api from '../api.js';
import { showNotification } from '../notifications/notifications.js';
import { isLogged, loginProcedure, logoutProcedure, closeModal, openModal, restoreAuthState, initAuthListeners } from '../auth/auth.js';
import {
    predefinedTypes,
    emailIdMap,
    tgIdMap,
    defaultNotifyNames,
    disabledPredefinedTypes,
    saveDisabledTypes,
    getTableBody
} from '../utils/utils.js';
import {
    createMultiselect,
    initAllMultiselects,
    updateAllMultiselects,
    webhookUrls,
    webhookSelections,
    multiselectInstances,
    renderWebhookTable
} from '../webhooks/webhooks.js';

let isSaving = false;
let customRows = [];
let editingRowId = null;
let predefinedDescriptions = {};

export function hideDisabledRows() {
    disabledPredefinedTypes.forEach(type => {
        const row = document.querySelector(`tr[data-notify-type="${type}"]`);
        if (row) row.classList.add('hidden-element');
    });
}

function initExistingMultiselects() {
    const containers = document.querySelectorAll('.webhook-multiselect-container');
    multiselectInstances.length = 0;
    containers.forEach(container => {
        const notifyType = container.dataset.notifyType;
        const instance = createMultiselect(container, notifyType);
        multiselectInstances.push({ notifyType, instance });
    });
}

function updateWebhookUrlsFromExistingRows() {
    const containers = document.querySelectorAll('.webhook-multiselect-container');
    const urls = new Set();
    containers.forEach(container => {
        const selects = container.querySelectorAll('.option-item input[type="checkbox"]');
        selects.forEach(cb => {
            if (cb.checked) urls.add(cb.value);
        });
    });
    webhookUrls.length = 0;
    webhookUrls.push(...urls);
}

export async function GetSettings() {
    console.log("GetSettings called");
    try {
        const result = await api.getNotifySettings();
        console.log('Result:', result);

        const tbody = getTableBody();
        if (!tbody) return;

        if (!result.data || !Array.isArray(result.data) || result.data.length === 0) {
            console.log("No data or empty data, keeping hardcoded rows");
            initExistingMultiselects();
            updateWebhookUrlsFromExistingRows();
            renderWebhookTable();
            renderCustomRows();
            return;
        }

        const data = result.data;

        tbody.innerHTML = '';
        customRows = [];
        multiselectInstances.length = 0;

        data.forEach(item => {
            if (predefinedTypes.includes(item.notify_type)) {
                if (disabledPredefinedTypes.has(item.notify_type)) {
                    return;
                }

                if (!item.notify_description && defaultNotifyNames[item.notify_type]) {
                    item.notify_description = defaultNotifyNames[item.notify_type];
                }

                if (item.notify_description) {
                    predefinedDescriptions[item.notify_type] = item.notify_description;
                }

                const row = createPredefinedRow(item);
                tbody.appendChild(row);

                const container = row.querySelector('.webhook-multiselect-container');
                if (container) {
                    const nt = container.dataset.notifyType;
                    webhookSelections[nt] = item.webhook_urls || [];
                    const instance = createMultiselect(container, nt);
                    multiselectInstances.push({ notifyType: nt, instance });
                }
            } else {
                customRows.push({
                    id: Date.now() + Math.random(),
                    notify_type: item.notify_type,
                    description: item.notify_description || item.notify_type,
                    want_email: item.want_email,
                    want_telegram: item.want_telegram,
                    webhook_urls: item.webhook_urls || []
                });
            }
        });

        renderCustomRows();

        const allUrls = [];
        data.forEach(item => {
            (item.webhook_urls || []).forEach(url => {
                if (url && !allUrls.includes(url)) allUrls.push(url);
            });
        });
        webhookUrls.length = 0;
        webhookUrls.push(...allUrls);

        updateAllMultiselects();

        renderWebhookTable();

    } catch (error) {
        console.error('Error:', error);
        renderWebhookTable();
        renderCustomRows();
    }
}

export async function CompleteSetup() {
    console.log('CompleteSetup called', new Date().toISOString());
    if (isSaving) return;
    isSaving = true;

    if (!isLogged) {
        console.log("Необходимо войти в аккаунт!");
        showNotification('error', 'Ошибка', 'Необходимо войти в аккаунт!');
        isSaving = false;
        return;
    }

    let payload = {
        "data": []
    };

    const activeTypes = predefinedTypes.filter(type => !disabledPredefinedTypes.has(type));

    activeTypes.forEach(type => {
        const emailId = emailIdMap[type];
        const tgId = tgIdMap[type];
        const customDescription = predefinedDescriptions[type];

        payload.data.push({
            "notify_type": type,
            "want_email": document.getElementById(emailId) ? document.getElementById(emailId).checked : false,
            "want_telegram": document.getElementById(tgId) ? document.getElementById(tgId).checked : false,
            "webhook_urls": webhookSelections[type] || [],
            ...(customDescription ? { "notify_description": customDescription } : {}),
        });
    });

    customRows.forEach(row => {
        if (row.notify_type && row.description) {
            payload.data.push({
                "notify_type": row.notify_type,
                "notify_description": row.description,
                "want_email": row.want_email,
                "want_telegram": row.want_telegram,
                "webhook_urls": webhookSelections[row.notify_type] || []
            });
        }
    });

    console.log("Final payload:", JSON.stringify(payload, null, 2));

    try {
        const response = await api.saveNotifySettings(payload);

        if (!response.ok) {
            showNotification('error', 'Ошибка', 'При обращении к БД произошла ошибка');
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        showNotification('success', 'Успех', 'Настройки сохранены');
    } catch (error) {
        console.error('Error:', error);
        showNotification('error', 'Ошибка', 'Не удалось сохранить настройки');
    } finally {
        isSaving = false;
    }
}

export function addCustomRow() {
    const newRow = {
        id: Date.now() + Math.random(),
        notify_type: '',
        description: '',
        want_email: false,
        want_telegram: false,
        webhook_urls: [],
        isNew: true
    };

    customRows.push(newRow);
    console.log('Добавлена строка, customRows теперь:', customRows);
    editingRowId = newRow.id;
    renderCustomRows();

    setTimeout(() => {
        const firstInput = document.querySelector(`tr[data-id="${newRow.id}"] input`);
        if (firstInput) firstInput.focus();
    }, 100);
}

export function deleteCustomRow(id) {
    customRows = customRows.filter(row => row.id !== id);
    renderCustomRows();
}

function createPredefinedRow(item) {
    const tr = document.createElement('tr');
    tr.dataset.notifyType = item.notify_type;

    const desc = item.notify_description || defaultNotifyNames[item.notify_type] || item.notify_type;

    const nameTd = document.createElement('td');
    nameTd.className = 'table-cell table-cell-left';
    nameTd.innerHTML = `<p class="text">${desc}</p>`;

    const idTd = document.createElement('td');
    idTd.className = 'table-cell table-cell-left';
    idTd.textContent = item.notify_type;

    const tgTd = document.createElement('td');
    tgTd.className = 'table-cell';
    const tgChk = document.createElement('input');
    tgChk.type = 'checkbox';
    tgChk.id = tgIdMap[item.notify_type] || `tg_${item.notify_type}`;
    tgChk.checked = item.want_telegram;
    tgTd.appendChild(tgChk);

    const emailTd = document.createElement('td');
    emailTd.className = 'table-cell';
    const emailChk = document.createElement('input');
    emailChk.type = 'checkbox';
    emailChk.id = emailIdMap[item.notify_type] || `email_${item.notify_type}`;
    emailChk.checked = item.want_email;
    emailTd.appendChild(emailChk);

    const webhookTd = document.createElement('td');
    webhookTd.className = 'table-cell';
    const webhookContainer = document.createElement('div');
    webhookContainer.className = 'webhook-multiselect-container';
    webhookContainer.dataset.notifyType = item.notify_type;
    webhookTd.appendChild(webhookContainer);

    const actionsTd = document.createElement('td');
    actionsTd.className = 'table-cell table-cell-center';
    const actionsDiv = document.createElement('div');
    actionsDiv.className = 'actions-container';
    actionsDiv.innerHTML = `
        <button class="webhook-action-btn edit" data-notify-type="${item.notify_type}" title="Редактировать">✎</button>
        <button class="webhook-action-btn delete delete-row-btn" data-notify-type="${item.notify_type}" title="Удалить">×</button>
    `;
    actionsTd.appendChild(actionsDiv);

    tr.appendChild(nameTd);
    tr.appendChild(idTd);
    tr.appendChild(tgTd);
    tr.appendChild(emailTd);
    tr.appendChild(webhookTd);
    tr.appendChild(actionsTd);

    tr._nameCell = nameTd;
    tr._actionsContainer = actionsDiv;

    return tr;
}

export function editPredefinedRow(notifyType) {
    const row = document.querySelector(`tr[data-notify-type="${notifyType}"]`);
    if (!row) return;

    const textCell = row._nameCell || row.querySelector('td:first-child');
    const textElem = textCell.querySelector('.text');
    const originalText = textElem.textContent;

    const input = document.createElement('input');
    input.type = 'text';
    input.value = originalText;
    input.className = 'custom-row-input custom-row-input-full';
    textCell.innerHTML = '';
    textCell.appendChild(input);
    input.focus();

    const actionsCell = row._actionsContainer || row.querySelector('.actions-container');
    actionsCell.innerHTML = '';

    const saveBtn = document.createElement('button');
    saveBtn.className = 'webhook-action-btn save';
    saveBtn.textContent = '✓';
    saveBtn.title = 'Сохранить изменения';

    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'webhook-action-btn delete';
    cancelBtn.textContent = '✕';
    cancelBtn.title = 'Отмена';

    saveBtn.onclick = function() {
        const newText = input.value.trim();
        if (newText) {
            textCell.innerHTML = `<p class="text">${newText}</p>`;
            predefinedDescriptions[notifyType] = newText;
            restorePredefinedButtons(notifyType);
            showNotification('success', 'Сохранено', 'Название параметра обновлено');
        } else {
            showNotification('error', 'Ошибка', 'Название не может быть пустым');
        }
    };

    cancelBtn.onclick = function() {
        textCell.innerHTML = `<p class="text">${originalText}</p>`;
        restorePredefinedButtons(notifyType);
    };

    actionsCell.appendChild(saveBtn);
    actionsCell.appendChild(cancelBtn);
}

function restorePredefinedButtons(notifyType) {
    const row = document.querySelector(`tr[data-notify-type="${notifyType}"]`);
    if (!row) return;
    const actionsCell = row._actionsContainer || row.querySelector('.actions-container');
    actionsCell.innerHTML = `
        <button class="webhook-action-btn edit" data-notify-type="${notifyType}" title="Редактировать">✎</button>
        <button class="webhook-action-btn delete delete-row-btn" data-notify-type="${notifyType}" title="Удалить">×</button>
    `;
    const newEditBtn = actionsCell.querySelector('.edit');
    const newDeleteBtn = actionsCell.querySelector('.delete-row-btn');
    if (newEditBtn) newEditBtn.onclick = () => editPredefinedRow(notifyType);
    if (newDeleteBtn) newDeleteBtn.onclick = () => deletePredefinedRow(notifyType);
}

export function deletePredefinedRow(notifyType) {
    if (confirm(`Вы уверены, что хотите удалить параметр ${notifyType}?`)) {
        disabledPredefinedTypes.add(notifyType);
        saveDisabledTypes();
        const row = document.querySelector(`tr[data-notify-type="${notifyType}"]`);
        if (row) row.remove();
        delete webhookSelections[notifyType];
        showNotification('success', 'Удалено', 'Параметр удалён');
    }
}

function updateCustomRow(id, field, value) {
    const row = customRows.find(r => r.id === id);
    if (row) {
        row[field] = value;
        if (field === 'notify_type' && !row.webhook_urls) {
            row.webhook_urls = [];
        }
    }
}

export function renderCustomRows() {
    const tbody = getTableBody();
    if (!tbody) return;

    const existingCustomRows = tbody.querySelectorAll('.custom-row');
    existingCustomRows.forEach(row => row.remove());

    customRows.forEach((row) => {
        const tr = document.createElement('tr');
        tr.className = 'custom-row';
        tr.dataset.id = row.id;

        const nameCell = document.createElement('td');
        nameCell.className = 'table-cell table-cell-left';

        if (row.isNew || editingRowId === row.id) {
            const nameInput = document.createElement('input');
            nameInput.type = 'text';
            nameInput.className = 'custom-row-input';
            nameInput.placeholder = 'Введите название';
            nameInput.value = row.description || '';
            nameInput.addEventListener('input', (e) => {
                updateCustomRow(row.id, 'description', e.target.value);
            });
            nameCell.appendChild(nameInput);
        } else {
            const nameText = document.createElement('span');
            nameText.textContent = row.description || '';
            nameCell.appendChild(nameText);
        }

        const idCell = document.createElement('td');
        idCell.className = 'table-cell table-cell-left';

        if (row.isNew || editingRowId === row.id) {
            const idInput = document.createElement('input');
            idInput.type = 'text';
            idInput.className = 'custom-row-input';
            idInput.placeholder = 'Введите ID';
            idInput.value = row.notify_type || '';
            idInput.addEventListener('input', (e) => {
                updateCustomRow(row.id, 'notify_type', e.target.value);
            });
            idCell.appendChild(idInput);
        } else {
            const idText = document.createElement('span');
            idText.textContent = row.notify_type || '';
            idCell.appendChild(idText);
        }

        const tgCell = document.createElement('td');
        tgCell.className = 'table-cell';
        const tgCheckbox = document.createElement('input');
        tgCheckbox.type = 'checkbox';
        tgCheckbox.checked = row.want_telegram;
        tgCheckbox.addEventListener('change', (e) => {
            updateCustomRow(row.id, 'want_telegram', e.target.checked);
        });
        tgCell.appendChild(tgCheckbox);

        const emailCell = document.createElement('td');
        emailCell.className = 'table-cell';
        const emailCheckbox = document.createElement('input');
        emailCheckbox.type = 'checkbox';
        emailCheckbox.checked = row.want_email;
        emailCheckbox.addEventListener('change', (e) => {
            updateCustomRow(row.id, 'want_email', e.target.checked);
        });
        emailCell.appendChild(emailCheckbox);

        const webhookCell = document.createElement('td');
        webhookCell.className = 'table-cell';
        const webhookContainer = document.createElement('div');
        webhookContainer.className = 'webhook-multiselect-container webhook-multiselect-min-width';
        webhookContainer.dataset.notifyType = row.notify_type || `custom_${row.id}`;
        webhookCell.appendChild(webhookContainer);

        const actionsCell = document.createElement('td');
        actionsCell.className = 'table-cell table-cell-center';

        const actionsContainer = document.createElement('div');
        actionsContainer.className = 'actions-container';

        if (row.isNew || editingRowId === row.id) {
            const saveBtn = document.createElement('button');
            saveBtn.className = 'webhook-action-btn save webhook-action-btn-margin-left';
            saveBtn.textContent = '✓';
            saveBtn.title = row.isNew ? 'Добавить параметр' : 'Сохранить изменения';

            saveBtn.addEventListener('click', () => {
                console.log('Сохранение, customRows до:', customRows);
                if (row.isNew) {
                    if (!row.notify_type || !row.description) {
                        showNotification('error', 'Ошибка', 'Заполните все поля');
                        return;
                    }
                    if (predefinedTypes.includes(row.notify_type) ||
                        customRows.some(r => r.id !== row.id && r.notify_type === row.notify_type)) {
                        showNotification('error', 'Ошибка', 'Параметр с таким ID уже существует');
                        return;
                    }
                    delete row.isNew;
                    editingRowId = null;
                    showNotification('success', 'Добавлено', 'Параметр добавлен');
                } else {
                    editingRowId = null;
                    showNotification('success', 'Сохранено', 'Изменения сохранены');
                }
                renderCustomRows();
                console.log('customRows после сохранения:', customRows);
            });

            actionsContainer.appendChild(saveBtn);
        } else {
            const editBtn = document.createElement('button');
            editBtn.className = 'webhook-action-btn edit webhook-action-btn-margin-right';
            editBtn.textContent = '✎';
            editBtn.title = 'Редактировать';
            editBtn.addEventListener('click', () => {
                editingRowId = row.id;
                renderCustomRows();
            });

            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'webhook-action-btn delete';
            deleteBtn.textContent = '✕';
            deleteBtn.title = 'Удалить';
            deleteBtn.addEventListener('click', () => {
                if (confirm(`Вы уверены, что хотите удалить параметр "${row.description}"?`)) {
                    deleteCustomRow(row.id);
                    showNotification('success', 'Удалено', 'Параметр удален');
                }
            });

            actionsContainer.appendChild(editBtn);
            actionsContainer.appendChild(deleteBtn);
        }

        actionsCell.appendChild(actionsContainer);

        tr.appendChild(nameCell);
        tr.appendChild(idCell);
        tr.appendChild(tgCell);
        tr.appendChild(emailCell);
        tr.appendChild(webhookCell);
        tr.appendChild(actionsCell);

        tbody.appendChild(tr);

        setTimeout(() => {
            if (!webhookSelections[webhookContainer.dataset.notifyType]) {
                webhookSelections[webhookContainer.dataset.notifyType] = row.webhook_urls || [];
            }
            const instance = createMultiselect(webhookContainer, webhookContainer.dataset.notifyType);
            multiselectInstances.push({
                notifyType: webhookContainer.dataset.notifyType,
                instance
            });
        }, 0);
    });

    setTimeout(() => {
        const editButtons = document.querySelectorAll('.webhook-action-btn.edit[data-notify-type]');
        editButtons.forEach(button => {
            if (!button.hasAttribute('data-handler-added')) {
                button.setAttribute('data-handler-added', 'true');
                button.addEventListener('click', function() {
                    const notifyType = this.dataset.notifyType;
                    editPredefinedRow(notifyType);
                });
            }
        });

        const deleteButtons = document.querySelectorAll('.delete-row-btn');
        deleteButtons.forEach(button => {
            button.addEventListener('click', function() {
                const notifyType = this.dataset.notifyType;
                if (notifyType && confirm(`Вы уверены, что хотите удалить параметр ${notifyType}?`)) {
                    deletePredefinedRow(notifyType);
                }
            });
        });
    }, 100);
}

export function initSettingsListeners() {
    const addCustomRowBtn = document.getElementById('addCustomRowBtn');
    if (addCustomRowBtn) {
        addCustomRowBtn.addEventListener('click', addCustomRow);
    }

    const saveBtn = document.getElementById('saveSettingsBtn');
    if (saveBtn) saveBtn.addEventListener('click', CompleteSetup);
}