import './webhooks.css';
import { showNotification } from '../notifications/notifications.js';

/**
 * Массив всех доступных вебхуков
 * @type {Array<string>}
 */
export let webhookUrls = [];

/**
 * Объект, хранящий выбранные вебхуки для каждого типа уведомления
 * Ключ - тип уведомления, значение - массив выбранных URL
 * @type {Object<string, Array<string>>}
 */
export let webhookSelections = {};

/**
 * Массив экземпляров мультиселектов для каждого типа уведомления
 * @type {Array<{notifyType: string, instance: Object}>}
 */
export let multiselectInstances = [];

/**
 * Индекс редактируемого вебхука в таблице (null если ничего не редактируется)
 * @type {number|null}
 */
let editingIndex = null;

/**
 * Создает мультиселект для выбора вебхуков для конкретного типа уведомления
 * @param {HTMLElement} container - Элемент контейнера, в котором будет создан мультиселект
 * @param {string} notifyType - Тип уведомления, для которого создается мультиселект
 * @returns {Object} - Объект с методами для управления мультиселектом (renderOptions, updateDisplayText, getSelected, setSelected)
 */
export function createMultiselect(container, notifyType) {
    if (!webhookSelections[notifyType]) {
        webhookSelections[notifyType] = [];
    }

    const wrapper = document.createElement('div');
    wrapper.className = 'custom-select';

    const display = document.createElement('div');
    display.className = 'select-display';
    const placeholder = document.createElement('span');
    placeholder.className = 'select-placeholder';
    placeholder.textContent = 'Выберите URL';
    const arrow = document.createElement('span');
    arrow.className = 'select-arrow';
    arrow.textContent = '▾';
    display.appendChild(placeholder);
    display.appendChild(arrow);

    const optionsPanel = document.createElement('div');
    optionsPanel.className = 'select-options';

    function renderOptions() {
        optionsPanel.innerHTML = '';
        if (webhookUrls.length === 0) {
            const emptyMsg = document.createElement('div');
            emptyMsg.className = 'option-item option-item-empty';
            emptyMsg.textContent = 'Нет доступных URL';
            optionsPanel.appendChild(emptyMsg);
            return;
        }
        webhookUrls.forEach(url => {
            const label = document.createElement('label');
            label.className = 'option-item';
            const cb = document.createElement('input');
            cb.type = 'checkbox';
            cb.value = url;
            cb.checked = webhookSelections[notifyType].includes(url);
            label.appendChild(cb);
            label.appendChild(document.createTextNode(url));

            label.addEventListener('click', function(e) {
                e.preventDefault();
                const checkbox = this.querySelector('input[type="checkbox"]');
                if (checkbox) {
                    checkbox.checked = !checkbox.checked;
                    const selected = webhookSelections[notifyType];
                    const url = checkbox.value;
                    if (checkbox.checked) {
                        if (!selected.includes(url)) {
                            selected.push(url);
                        }
                    } else {
                        const idx = selected.indexOf(url);
                        if (idx !== -1) selected.splice(idx, 1);
                    }
                    updateDisplayText();
                }
            });

            optionsPanel.appendChild(label);

            cb.addEventListener('change', function(e) {
                e.stopPropagation();
                const selected = webhookSelections[notifyType];
                if (this.checked) {
                    if (!selected.includes(this.value)) {
                        selected.push(this.value);
                    }
                } else {
                    const idx = selected.indexOf(this.value);
                    if (idx !== -1) selected.splice(idx, 1);
                }
                updateDisplayText();
            });
        });
    }

    function updateDisplayText() {
        const selected = webhookSelections[notifyType] || [];
        if (selected.length === 0) {
            placeholder.textContent = 'Выберите URL';
            placeholder.classList.remove('has-value');
        } else {
            placeholder.textContent = selected.join(', ');
            placeholder.classList.add('has-value');
        }
    }

    function togglePanel(forceState) {
        if (typeof forceState === 'boolean') {
            wrapper.classList.toggle('open', forceState);
        } else {
            wrapper.classList.toggle('open');
        }
    }

    display.addEventListener('click', function(e) {
        e.stopPropagation();
        togglePanel();
    });

    document.addEventListener('click', function(e) {
        if (!wrapper.contains(e.target)) {
            togglePanel(false);
        }
    });

    wrapper.appendChild(display);
    wrapper.appendChild(optionsPanel);
    container.innerHTML = '';
    container.appendChild(wrapper);

    renderOptions();
    updateDisplayText();

    return {
        renderOptions: renderOptions,
        updateDisplayText: updateDisplayText,
        getSelected: () => webhookSelections[notifyType] || [],
        setSelected: (urls) => {
            webhookSelections[notifyType] = urls || [];
            renderOptions();
            updateDisplayText();
        }
    };
}

/**
 * Инициализирует все мультиселекты на странице
 * Находит все элементы с классом 'webhook-multiselect-container',
 * создает для каждого экземпляр мультиселекта и сохраняет их в массиве multiselectInstances
 */
export function initAllMultiselects() {
    const containers = document.querySelectorAll('.webhook-multiselect-container');
    multiselectInstances = [];
    containers.forEach(container => {
        const notifyType = container.dataset.notifyType;
        const instance = createMultiselect(container, notifyType);
        multiselectInstances.push({ notifyType, instance });
    });
}

/**
 * Очищает выборы вебхуков, удаляя те, которые больше не доступны
 * После удаления вебхуков из общего списка, эта функция удаляет их из всех выборов
 * @see webhookUrls
 * @see webhookSelections
 */
export function cleanupSelections() {
    const currentUrls = new Set(webhookUrls);
    Object.keys(webhookSelections).forEach(notifyType => {
        webhookSelections[notifyType] = webhookSelections[notifyType].filter(url => currentUrls.has(url));
    });
}

/**
 * Обновляет все мультиселекты на странице
 * Сначала очищает выборы (удаляет недоступные URL), затем обновляет отображение каждого мультиселекта
 * @see cleanupSelections
 */
export function updateAllMultiselects() {
    cleanupSelections();
    multiselectInstances.forEach(({ instance }) => {
        instance.renderOptions();
        instance.updateDisplayText();
    });
}

export function renderWebhookTable() {
    const tbody = document.getElementById('webhookTableBody');
    if (!tbody) {
        console.warn('webhookTableBody не найден!');
        return;
    }

    const rows = [];
    webhookUrls.forEach((url, index) => {
        const isEditing = (editingIndex === index);
        rows.push({ url, index, isEditing, isNew: false });
    });

    if (editingIndex === null) {
        rows.push({ url: '', index: -1, isEditing: false, isNew: true });
    }

    tbody.innerHTML = '';
    rows.forEach((row) => {
        const tr = document.createElement('tr');

        const td = document.createElement('td');
        const container = document.createElement('div');
        container.className = 'webhook-row';

        const input = document.createElement('input');
        input.type = 'text';
        input.className = 'webhook-url-input';
        input.placeholder = 'Введите ссылку';
        input.value = row.url;
        input.disabled = !row.isEditing && !row.isNew;
        container.appendChild(input);

        if (row.isNew && !row.isEditing) {
            const saveBtn = document.createElement('button');
            saveBtn.className = 'webhook-action-btn save hidden-element';
            saveBtn.textContent = '✓';
            saveBtn.title = 'Сохранить новый URL';
            saveBtn.classList.add('webhook-action-btn-margin-left');

            input.addEventListener('input', function() {
                if (this.value.trim()) {
                    saveBtn.classList.remove('hidden-element');
                } else {
                    saveBtn.classList.add('hidden-element');
                }
            });

            saveBtn.addEventListener('click', function(e) {
                e.stopPropagation();
                const val = input.value.trim();
                if (!val) return;

                if (webhookUrls.includes(val)) {
                    showNotification('error', 'Ошибка', 'Такой URL уже существует');
                    return;
                }

                webhookUrls.push(val);
                editingIndex = null;
                renderWebhookTable();
                showNotification('success', 'Добавлено', 'Новый вебхук добавлен');
            });

            container.appendChild(saveBtn);
        } else if (!row.isNew) {
            if (row.isEditing) {
                const saveEditBtn = document.createElement('button');
                saveEditBtn.className = 'webhook-action-btn save webhook-action-btn-margin-left';
                saveEditBtn.textContent = '✓';
                saveEditBtn.title = 'Сохранить изменения';
                saveEditBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const val = input.value.trim();
                    if (!val) {
                        showNotification('error', 'Ошибка', 'URL не может быть пустым');
                        return;
                    }

                    if (val !== webhookUrls[row.index] && webhookUrls.includes(val)) {
                        showNotification('error', 'Ошибка', 'Такой URL уже существует');
                        return;
                    }

                    webhookUrls[row.index] = val;
                    editingIndex = null;
                    renderWebhookTable();
                    showNotification('success', 'Изменено', 'URL обновлён');
                });
                container.appendChild(saveEditBtn);
            } else {
                const editBtn = document.createElement('button');
                editBtn.className = 'webhook-action-btn edit webhook-action-btn-margin-right';
                editBtn.textContent = '✎';
                editBtn.title = 'Редактировать URL';
                editBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    editingIndex = row.index;
                    renderWebhookTable();
                    setTimeout(() => {
                        const inp = document.querySelector(`#webhookTableBody tr[data-index="${row.index}"] .webhook-url-input`);
                        if (inp) inp.focus();
                    }, 50);
                });

                const deleteBtn = document.createElement('button');
                deleteBtn.className = 'webhook-action-btn delete';
                deleteBtn.textContent = '✕';
                deleteBtn.title = 'Удалить URL';
                deleteBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const index = row.index;
                    if (confirm(`Удалить вебхук "${webhookUrls[index]}"?`)) {
                        webhookUrls.splice(index, 1);
                        if (editingIndex === index) editingIndex = null;
                        else if (editingIndex !== null && editingIndex > index) editingIndex--;
                        renderWebhookTable();
                        showNotification('success', 'Удалено', 'Вебхук удалён');
                    }
                });

                container.appendChild(editBtn);
                container.appendChild(deleteBtn);
            }
        }

        td.appendChild(container);
        tr.appendChild(td);
        tr.dataset.index = row.index;
        tbody.appendChild(tr);
    });
    updateAllMultiselects();
}