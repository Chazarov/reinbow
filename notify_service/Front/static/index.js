import './styles/main.css';
import './styles/layout.css';
import './styles/components.css';

let isSaving = false;
let customRows = [];
let editingRowId = null;

// Массив предопределенных типов для отслеживания
const predefinedTypes = ["user_register", "user_login", "admin_newImg", "user_imgVerdict"];

(function initNotifications() {    
    const HEIGHT = 80;             
    const GAP = 20;                
    const LIFETIME = 3000;         
    const CONTAINER_WIDTH = '18%';
    const LEFT_OFFSET = '2%';
    const BOTTOM_OFFSET = 50;       

    const template = document.getElementById('notification');
    if (template) {
        template.style.display = 'none';
    }

    const notifications = []; 

    function render() {
        notifications.forEach((item, index) => {
            const bottomPos = BOTTOM_OFFSET + index * (HEIGHT + GAP);
            item.element.style.bottom = bottomPos + 'px';
            item.element.style.opacity = '1';
        });
    }

    function removeOldest() {
        if (notifications.length === 0) return;

        const oldest = notifications.pop();
        const el = oldest.element;

        el.style.transition = 'bottom 0.5s ease, opacity 0.5s ease';
        el.style.bottom = (window.innerHeight + HEIGHT) + 'px';
        el.style.opacity = '0';

        const onFinish = () => {
            el.remove();
            el.removeEventListener('transitionend', onFinish);
        };
        el.addEventListener('transitionend', onFinish);

        render(); 
    }

    window.showNotification = function(type, title, description) {
        if (!template) return;

        const clone = template.cloneNode(true);
        clone.id = '';
        clone.style.display = 'block';
        clone.style.position = 'fixed';
        clone.style.left = LEFT_OFFSET;
        clone.style.width = 'auto';
        clone.style.maxWidth = '500px';
        clone.style.minWidth = '200px';
        clone.style.padding = '12px 20px';
        clone.style.borderRadius = '5px';
        clone.style.overflow = 'hidden';
        clone.style.zIndex = '2';
        clone.style.pointerEvents = 'none';

        clone.style.backgroundColor = type === 'success'
            ? 'rgba(0, 255, 0, 0.75)'
            : 'rgba(255, 0, 0, 0.75)';

        const titleSpan = clone.querySelector('.errorText');
        const descSpan = clone.querySelector('.errDesc');
        if (titleSpan) titleSpan.textContent = title;
        if (descSpan) descSpan.textContent = description;

        clone.style.transition = 'none';
        clone.style.bottom = -(HEIGHT + 100) + 'px';
        clone.style.opacity = '0';

        notifications.unshift({ element: clone });

        document.body.appendChild(clone);

        clone.offsetHeight;

        clone.style.transition = 'bottom 0.5s ease, opacity 0.5s ease';
        render();

        setTimeout(() => {
            removeOldest();
        }, LIFETIME);
    };
})();

// Функция для получения настроек из базы
async function GetSettings() {
    console.log("GetSettings called");
    try {
        const response = await fetchWithAuth("/api/get_notify_settings");
        const result = await response.json();
        console.log('Result:', result);

        // Очищаем пользовательские строки
        customRows = [];

        if (!result.data || !Array.isArray(result.data)) {
            console.log("No data or invalid data format");
            renderWebhookTable();
            renderCustomRows();
            return;
        }

        const data = result.data;

        // Обрабатываем предопределенные типы
        data.forEach((item) => {
            if (predefinedTypes.includes(item.notify_type)) {
                const notify_type = item.notify_type;
                const want_email = item.want_email;
                const want_telegram = item.want_telegram;

                switch (notify_type) {
                    case "user_register":
                        document.getElementById('email_reg').checked = want_email;
                        document.getElementById('tg_reg').checked = want_telegram;
                        break;
                    case "user_login":
                        document.getElementById('email_login').checked = want_email;
                        document.getElementById('tg_login').checked = want_telegram;
                        break;
                    case "admin_newImg":
                        document.getElementById('email_admin_newImg').checked = want_email;
                        document.getElementById('tg_admin_newImg').checked = want_telegram;
                        break;
                    case "user_imgVerdict":
                        document.getElementById('email_user_imgVerdict').checked = want_email;
                        document.getElementById('tg_user_imgVerdict').checked = want_telegram;
                        break;
                }
            } else {
                // Добавляем как пользовательскую строку
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

        // Собираем все уникальные URL
        const allUrls = [];
        data.forEach(item => {
            const urls = item.webhook_urls || [];
            if (Array.isArray(urls)) {
                urls.forEach(url => {
                    if (url && !allUrls.includes(url)) allUrls.push(url);
                });
                webhookSelections[item.notify_type] = urls;
            } else {
                webhookSelections[item.notify_type] = [];
            }
        });
        webhookUrls = allUrls;

        // Пересоздаём мультиселекты
        multiselectInstances = [];
        initAllMultiselects();

        // Рендерим таблицу URL и пользовательские строки
        renderWebhookTable();
        renderCustomRows();

        // Обновляем отображение мультиселектов
        updateAllMultiselects();

    } catch (error) {
        console.error('Error:', error);
        renderWebhookTable();
        renderCustomRows();
    }
}

let webhookSelections = {};

// Функция для сохранения настроек
async function CompleteSetup() {
    console.log('CompleteSetup called', new Date().toISOString());
    if (isSaving) return;
    isSaving = true;

    if (!isLogged) {
        console.log("Необходимо войти в аккаунт!");
        showNotification('error', 'Ошибка', 'Необходимо войти в аккаунт!');
        isSaving = false;
        return;
    }

    // Базовые данные для предопределенных типов
    let payload = {
        "data": []
    };

    // Добавляем предопределенные типы
    predefinedTypes.forEach(type => {
        let emailId, tgId;
        switch(type) {
            case "user_register":
                emailId = 'email_reg';
                tgId = 'tg_reg';
                break;
            case "user_login":
                emailId = 'email_login';
                tgId = 'tg_login';
                break;
            case "admin_newImg":
                emailId = 'email_admin_newImg';
                tgId = 'tg_admin_newImg';
                break;
            case "user_imgVerdict":
                emailId = 'email_user_imgVerdict';
                tgId = 'tg_user_imgVerdict';
                break;
        }

        payload.data.push({
            "notify_type": type,
            "want_email": document.getElementById(emailId).checked,
            "want_telegram": document.getElementById(tgId).checked,
            "webhook_urls": webhookSelections[type] || []
        });
    });

    // Добавляем пользовательские строки
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
        const response = await fetchWithAuth("/api/notify_types", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

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

// Удаляет из webhookSelections URL, которых нет в webhookUrls
function cleanupSelections() {
    const currentUrls = new Set(webhookUrls);
    Object.keys(webhookSelections).forEach(notifyType => {
        webhookSelections[notifyType] = webhookSelections[notifyType].filter(url => currentUrls.has(url));
    });
}

function createMultiselect(container, notifyType) {
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
            emptyMsg.className = 'option-item';
            emptyMsg.textContent = 'Нет доступных URL';
            emptyMsg.style.color = '#999';
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
            placeholder.style.color = '#999';
        } else {
            placeholder.textContent = selected.join(', ');
            placeholder.style.color = '#333';
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

// Глобальный массив для хранения объектов мультиселектов
let multiselectInstances = [];

function initAllMultiselects() {
    const containers = document.querySelectorAll('.webhook-multiselect-container');
    multiselectInstances = [];
    containers.forEach(container => {
        const notifyType = container.dataset.notifyType;
        const instance = createMultiselect(container, notifyType);
        multiselectInstances.push({ notifyType, instance });
    });
}

// Функция обновления всех мультиселектов при изменении списка URL
function updateAllMultiselects() {
    cleanupSelections(); // удаляем устаревшие URL из выбранных
    multiselectInstances.forEach(({ instance }) => {
        instance.renderOptions();
        instance.updateDisplayText();
    });
}

// Функция для анимации закрытия всплывающего окна входа в аккаунт
function closeModal() {
    let modal_window = document.getElementById('modalWindow');
    let black_background = document.getElementById('blackBackground');
    modal_window.classList.remove('open');
    modal_window.classList.add('closed');
    black_background.classList.remove('open');
    black_background.classList.add('closed');
    // Сброс инлайновых стилей (на случай, если они остались)
    modal_window.style.zIndex = '';
    modal_window.style.opacity = '';
    black_background.style.pointerEvents = '';
    black_background.style.background = '';
    black_background.style.backdropFilter = '';
}

// Функция для анимации открытия всплывающего окна входа в аккаунт
function openModal() {
    let modal_window = document.getElementById('modalWindow');
    let black_background = document.getElementById('blackBackground');
    modal_window.classList.remove('closed');
    modal_window.classList.add('open');
    black_background.classList.remove('closed');
    black_background.classList.add('open');
    // Сброс инлайновых стилей, чтобы классы работали
    modal_window.style.zIndex = '';
    modal_window.style.opacity = '';
    black_background.style.pointerEvents = '';
    black_background.style.background = '';
    black_background.style.backdropFilter = '';
}

// Функция для процедуры входа в аккаунт
async function loginProcedure() {
    let login_value = document.getElementById('login_field').value.trim();
    let password_value = document.getElementById('password_field').value.trim();

    const payload = {
        "login": login_value,
        "password": password_value
    }

    const response = await fetch("/api/moderator_login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    }
    )

    const result = await response.json();
    console.log(result)
    if (result.success == true) {
        console.log("Вход успешен!");
        localStorage.setItem('token', result.token)
        isLogged = true;
        document.getElementById('not_Logged').style.display = "none";
        document.getElementById('isAuthorized').style.display = "block";
        closeModal();
        showNotification('success', 'Успех!', 'Успешный вход');
        GetSettings();
    } else {
        console.log("Неправильные данные!");
        showNotification('error', 'Ошибка', result.Error_message);
    }
}

// Функция для процедуры выхода из аккаунта

async function logoutProcedure() {
    localStorage.removeItem('token');
    isLogged = false;

    document.getElementById('login_field').value = "";
    document.getElementById('password_field').value = "";

    document.getElementById('not_Logged').style.display = "flex";
    document.getElementById('isAuthorized').style.display = "none";

    // Сбрасываем данные
    webhookUrls = [];
    webhookSelections = {};
    multiselectInstances = [];
    // Пересоздаём пустые мультиселекты
    initAllMultiselects();
    renderWebhookTable(); // таблица станет пустой

    console.log("Выход успешен!");
    closeModal();
    showNotification('success', 'Успех', 'Выход из аккаунта успешен!');
}

// Функция для добавления хедера Authorization
function fetchWithAuth(url, options={}) {
    const token = localStorage.getItem('token')

    if (token) {
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`
        }
    }

    return fetch(url, options)
}

function restoreAuthState() {
    const token = localStorage.getItem('token')
    
    if (token) {
        isLogged = true
        document.getElementById('not_Logged').style.display = "none";
        document.getElementById('isAuthorized').style.display = "block";
        console.log("Добро пожаловать, admin")
        initAllMultiselects();
        GetSettings()
    } else {
        isLogged = false
        document.getElementById('not_Logged').style.display = "flex";
        document.getElementById('isAuthorized').style.display = "none";
    }
}

// ========== Управление списком вебхуков ==========
let webhookUrls = [];          // массив сохранённых URL
let editingIndex = null;      // индекс редактируемой строки

function renderWebhookTable() {
    const tbody = document.getElementById('webhookTableBody');
    if (!tbody) {
        console.warn('webhookTableBody не найден!');
        return;
    }

    // Строим список строк: все сохранённые + одна пустая (если не редактируется существующая)
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

        // Поле ввода
        const input = document.createElement('input');
        input.type = 'text';
        input.className = 'webhook-url-input';
        input.placeholder = 'Введите ссылку';
        input.value = row.url;
        input.disabled = !row.isEditing && !row.isNew;
        container.appendChild(input);

        // ---- Кнопки ----
        if (row.isNew && !row.isEditing) {
            // Новая пустая строка
            const saveBtn = document.createElement('button');
            saveBtn.className = 'webhook-action-btn save';
            saveBtn.textContent = '✓';
            saveBtn.style.display = 'none';
            saveBtn.title = 'Сохранить новый URL';
            saveBtn.style.marginLeft = '5px';

            input.addEventListener('input', function() {
                saveBtn.style.display = this.value.trim() ? 'inline-flex' : 'none';
            });

            saveBtn.addEventListener('click', function(e) {
                e.stopPropagation();
                const val = input.value.trim();
                if (!val) return;
                webhookUrls.push(val);
                editingIndex = null;
                renderWebhookTable();
                showNotification('success', 'Добавлено', 'Новый вебхук добавлен');
            });

            container.appendChild(saveBtn);
        } else if (!row.isNew) {
            // Существующий URL
            if (row.isEditing) {
                // Режим редактирования – кнопка "сохранить изменения"
                const saveEditBtn = document.createElement('button');
                saveEditBtn.className = 'webhook-action-btn save';
                saveEditBtn.textContent = '✓';
                saveEditBtn.title = 'Сохранить изменения';
                saveEditBtn.style.marginLeft = '5px';
                saveEditBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const val = input.value.trim();
                    if (!val) {
                        showNotification('error', 'Ошибка', 'URL не может быть пустым');
                        return;
                    }
                    webhookUrls[row.index] = val;
                    editingIndex = null;
                    renderWebhookTable();
                    showNotification('success', 'Изменено', 'URL обновлён');
                });
                container.appendChild(saveEditBtn);
            } else {
                // Обычный режим – редактировать и удалить
                const editBtn = document.createElement('button');
                editBtn.className = 'webhook-action-btn edit';
                editBtn.textContent = '✎';
                editBtn.title = 'Редактировать URL';
                editBtn.style.marginRight = '5px';
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
        // Сохраняем индекс строки для поиска при фокусе
        tr.dataset.index = row.index;
        tbody.appendChild(tr);
    });
    updateAllMultiselects();
}

// Функция для добавления пользовательской строки
function addCustomRow() {
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
    
    // Фокус на первое поле ввода
    setTimeout(() => {
        const firstInput = document.querySelector(`tr[data-id="${newRow.id}"] input`);
        if (firstInput) firstInput.focus();
    }, 100);
}

// Функция для удаления пользовательской строки
function deleteCustomRow(id) {
    customRows = customRows.filter(row => row.id !== id);
    renderCustomRows();
}

// Функция для удаления предопределенного параметра
function deletePredefinedRow(notifyType) {
    customRows = customRows.filter(row => row.notify_type !== notifyType);
    delete webhookSelections[notifyType];
    showNotification('success', 'Удалено', 'Параметр удален');
    renderCustomRows();
}

// Функция для обновления значения в пользовательской строке (без перерисовки)
function updateCustomRow(id, field, value) {
    const row = customRows.find(r => r.id === id);
    if (row) {
        row[field] = value;
        if (field === 'notify_type' && !row.webhook_urls) {
            row.webhook_urls = [];
        }
    }
}

// Функция отрисовки пользовательских строк
function renderCustomRows() {
    const table = document.querySelector('.table_main');
    
    // Удаляем существующие пользовательские строки
    const existingCustomRows = document.querySelectorAll('.custom-row');
    existingCustomRows.forEach(row => row.remove());
    
    // Добавляем пользовательские строки
    customRows.forEach((row) => {
        const tr = document.createElement('tr');
        tr.className = 'custom-row';
        tr.dataset.id = row.id;
        
        // Ячейка с названием
        const nameCell = document.createElement('td');
        nameCell.style.border = '1px solid transparent';
        nameCell.style.padding = '8px';
        nameCell.style.textAlign = 'left';
        
        if (row.isNew || editingRowId === row.id) {
            // Режим редактирования - показываем поля ввода
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
            // Режим просмотра - показываем текст
            const nameText = document.createElement('span');
            nameText.textContent = row.description || '';
            nameCell.appendChild(nameText);
        }
        
        // Ячейка с ID параметра
        const idCell = document.createElement('td');
        idCell.style.border = '1px solid transparent';
        idCell.style.padding = '8px';
        idCell.style.textAlign = 'left';
        
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
        
        // Ячейка Telegram
        const tgCell = document.createElement('td');
        tgCell.style.border = '1px solid transparent';
        tgCell.style.padding = '8px';
        const tgCheckbox = document.createElement('input');
        tgCheckbox.type = 'checkbox';
        tgCheckbox.checked = row.want_telegram;
        tgCheckbox.addEventListener('change', (e) => {
            updateCustomRow(row.id, 'want_telegram', e.target.checked);
        });
        tgCell.appendChild(tgCheckbox);
        
        // Ячейка Email
        const emailCell = document.createElement('td');
        emailCell.style.border = '1px solid transparent';
        emailCell.style.padding = '8px';
        const emailCheckbox = document.createElement('input');
        emailCheckbox.type = 'checkbox';
        emailCheckbox.checked = row.want_email;
        emailCheckbox.addEventListener('change', (e) => {
            updateCustomRow(row.id, 'want_email', e.target.checked);
        });
        emailCell.appendChild(emailCheckbox);
        
        // Ячейка Webhook
        const webhookCell = document.createElement('td');
        webhookCell.style.border = '1px solid transparent';
        webhookCell.style.padding = '8px';
        const webhookContainer = document.createElement('div');
        webhookContainer.className = 'webhook-multiselect-container';
        webhookContainer.dataset.notifyType = row.notify_type || `custom_${row.id}`;
        webhookContainer.style.minWidth = '200px';
        webhookCell.appendChild(webhookContainer);
        
        // Ячейка с кнопками действий
        const actionsCell = document.createElement('td');
        actionsCell.style.border = '1px solid transparent';
        actionsCell.style.padding = '8px';
        actionsCell.style.textAlign = 'center';
        
        const actionsContainer = document.createElement('div');
        actionsContainer.className = 'actions-container';
        
        if (row.isNew || editingRowId === row.id) {
            // Кнопка сохранения
            const saveBtn = document.createElement('button');
            saveBtn.className = 'webhook-action-btn save';
            saveBtn.textContent = '✓';
            saveBtn.title = row.isNew ? 'Добавить параметр' : 'Сохранить изменения';
            saveBtn.style.marginLeft = '5px';
            
            saveBtn.addEventListener('click', () => {
                console.log('Сохранение, customRows до:', customRows);
                if (row.isNew) {
                    if (!row.notify_type || !row.description) {
                        showNotification('error', 'Ошибка', 'Заполните все поля');
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
            // Кнопки редактирования и удаления
            const editBtn = document.createElement('button');
            editBtn.className = 'webhook-action-btn edit';
            editBtn.textContent = '✎';
            editBtn.title = 'Редактировать';
            editBtn.style.marginRight = '5px';
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
        
        table.appendChild(tr);
        
        // Создаем мультиселект для этой строки
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
    
    // Добавляем обработчики для предопределенных строк
    setTimeout(() => {
        const editButtons = document.querySelectorAll('.webhook-action-btn.edit[data-notify-type]');
        editButtons.forEach(button => {
            // Проверяем, не добавлен ли уже обработчик
            if (!button.hasAttribute('data-handler-added')) {
                button.setAttribute('data-handler-added', 'true');
                button.addEventListener('click', function() {
                    const notifyType = this.dataset.notifyType;
                    showNotification('info', 'Информация', `Редактирование параметра ${notifyType} пока не реализовано для предопределенных строк`);
                });
            }
        });
    }, 100);
}

let isLogged = false;
restoreAuthState();

// Добавляем обработчики событий (гарантируем однократное выполнение)
if (!window._listenersAdded) {
    window._listenersAdded = true;

    document.addEventListener('DOMContentLoaded', function() {
        const addCustomRowBtn = document.getElementById('addCustomRowBtn');
        if (addCustomRowBtn) {
            addCustomRowBtn.addEventListener('click', addCustomRow);
        }
        
        // Добавляем обработчики для кнопок удаления предопределенных строк
        const deleteButtons = document.querySelectorAll('.delete-row-btn');
        deleteButtons.forEach(button => {
            button.addEventListener('click', function() {
                const notifyType = this.dataset.notifyType;
                if (notifyType && confirm(`Вы уверены, что хотите удалить параметр ${notifyType}?`)) {
                    deletePredefinedRow(notifyType);
                }
            });
        });
    });

    const openModalBtn = document.getElementById('openModalBtn');
    if (openModalBtn) openModalBtn.addEventListener('click', openModal);

    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) logoutBtn.addEventListener('click', logoutProcedure);

    const closeModalBtn = document.getElementById('closeModalBtn');
    if (closeModalBtn) closeModalBtn.addEventListener('click', closeModal);

    const loginBtn = document.getElementById('loginBtn');
    if (loginBtn) loginBtn.addEventListener('click', loginProcedure);

    const saveBtn = document.getElementById('saveSettingsBtn');
    if (saveBtn) saveBtn.addEventListener('click', CompleteSetup);
}

// Для просмотра из консоли
window.customRows = customRows;
window.predefinedTypes = predefinedTypes;
window.webhookSelections = webhookSelections;