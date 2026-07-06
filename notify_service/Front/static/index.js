import './styles/main.css';
import './styles/layout.css';
import './styles/components.css';
let isSaving = false;

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
        clone.style.width = CONTAINER_WIDTH;
        clone.style.height = HEIGHT + 'px';
        clone.style.margin = '0';
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

        // Проверяем, что data существует и это массив
        if (!result.data || !Array.isArray(result.data)) {
            console.log("No data or invalid data format");
            renderWebhookTable();
            return;
        }

        const data = result.data;

        // Восстанавливаем чекбоксы Telegram и Email
        data.forEach((item) => {
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
        });

        // 1. Собираем все уникальные URL и восстанавливаем webhookSelections
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

        // 2. Пересоздаём мультиселекты
        multiselectInstances = [];
        initAllMultiselects();

        // 3. Рендерим таблицу URL
        renderWebhookTable();

        // 4. Обновляем отображение мультиселектов
        updateAllMultiselects();

    } catch (error) {
        console.error('Error:', error);
        renderWebhookTable();
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

    // ... после проверки isLogged

    let payload = {
        "data": [
            {
                "notify_type": "user_register",
                "want_email": document.getElementById('email_reg').checked,
                "want_telegram": document.getElementById('tg_reg').checked,
                "webhook_urls": webhookSelections['user_register'] || []
            },
            {
                "notify_type": "user_login",
                "want_email": document.getElementById('email_login').checked,
                "want_telegram": document.getElementById('tg_login').checked,
                "webhook_urls": webhookSelections['user_login'] || []
            },
            {
                "notify_type": "admin_newImg",
                "want_email": document.getElementById('email_admin_newImg').checked,
                "want_telegram": document.getElementById('tg_admin_newImg').checked,
                "webhook_urls": webhookSelections['admin_newImg'] || []
            },
            {
                "notify_type": "user_imgVerdict",
                "want_email": document.getElementById('email_user_imgVerdict').checked,
                "want_telegram": document.getElementById('tg_user_imgVerdict').checked,
                "webhook_urls": webhookSelections['user_imgVerdict'] || []
            }
        ]
    };

    try {

        console.log("payload:", payload)

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
                // Предотвращаем двойное срабатывание (чтобы не вызвать обработчик дважды)
                e.preventDefault();
                const checkbox = this.querySelector('input[type="checkbox"]');
                if (checkbox) {
                    // Переключаем состояние чекбокса
                    checkbox.checked = !checkbox.checked;
                    // Вручную обновляем webhookSelections
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
                    // Обновляем отображение (placeholder)
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

let isLogged = false;
restoreAuthState();

// Добавляем обработчики событий (гарантируем однократное выполнение)
if (!window._listenersAdded) {
    window._listenersAdded = true;

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