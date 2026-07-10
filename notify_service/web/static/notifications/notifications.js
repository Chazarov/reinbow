import './notifications.css';

const HEIGHT = 80;
const GAP = 20;
const LIFETIME = 3000;
const BOTTOM_OFFSET = 50;

const template = document.getElementById('notification');
if (template) {
    template.classList.add('hidden');
}

const notifications = [];

/**
 * Обновляет позиции всех уведомлений на экране
 * Располагает уведомления вертикально с учетом отступов и высоты
 */
/**
 * Обновляет позиции всех уведомлений на экране
 * Располагает уведомления вертикально с учетом отступов и высоты
 */
function render() {
    notifications.forEach((item, index) => {
        item.element.classList.remove('out');
        for (let i = 0; i < 10; i++) {
            item.element.classList.remove('pos-' + i);
        }
        item.element.classList.add('pos-' + index, 'visible');
    });
}

/**
 * Удаляет самое старое уведомление из очереди
 * Применяет анимацию удаления и обновляет отображение оставшихся уведомлений
 */
function removeOldest() {
    if (notifications.length === 0) return;

    const oldest = notifications.pop();
    const el = oldest.element;

    el.classList.remove('visible');
    el.classList.add('out');

    const onFinish = () => {
        el.remove();
        el.removeEventListener('transitionend', onFinish);
    };
    el.addEventListener('transitionend', onFinish);

    render();
}

/**
 * Отображает уведомление на экране
 * @param {string} type - Тип уведомления (success, error, warning, info)
 * @param {string} title - Заголовок уведомления
 * @param {string} description - Описание уведомления
 */
export function showNotification(type, title, description) {
    if (!template) return;

    const clone = template.cloneNode(true);
    clone.id = '';
    clone.classList.add('notification', type);

    const titleSpan = clone.querySelector('.errorText');
    const descSpan = clone.querySelector('.errDesc');
    if (titleSpan) titleSpan.textContent = title;
    if (descSpan) descSpan.textContent = description;

    notifications.unshift({ element: clone });

    document.body.appendChild(clone);

    clone.classList.remove('hidden');
    clone.classList.add('visible');
    render();

    setTimeout(() => {
        removeOldest();
    }, LIFETIME);
}

window.showNotification = showNotification;