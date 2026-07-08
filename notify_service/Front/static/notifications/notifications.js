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

function render() {
    notifications.forEach((item, index) => {
        const bottomPos = BOTTOM_OFFSET + index * (HEIGHT + GAP);
        item.element.style.bottom = bottomPos + 'px';
        item.element.classList.add('visible');
    });
}

function removeOldest() {
    if (notifications.length === 0) return;

    const oldest = notifications.pop();
    const el = oldest.element;

    el.style.transition = 'bottom 0.5s ease, opacity 0.5s ease';
    el.style.bottom = (window.innerHeight + HEIGHT) + 'px';
    el.classList.remove('visible');

    const onFinish = () => {
        el.remove();
        el.removeEventListener('transitionend', onFinish);
    };
    el.addEventListener('transitionend', onFinish);

    render();
}

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