const headerStylesHref = new URL('./header.css', import.meta.url).href;

function ensureHeaderStyles() {
    if (document.querySelector('link[data-component="header"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = headerStylesHref;
    link.dataset.component = 'header';
    document.head.appendChild(link);
}

export function renderHeader(config = {}) {
    const { authorized = false } = config;
    ensureHeaderStyles();

    const header = document.createElement('header');
    header.className = 'site-header';

    const nav = document.createElement('nav');
    const title = document.createElement('a');
    title.href = '#/login';
    title.textContent = 'Image Service';

    const list = document.createElement('ul');
    const navLinks = [
        { href: '#/register', label: 'Регистрация' , authorized_only: false, guest_only: true},
        { href: '#/login', label: 'Вход', authorized_only: false, guest_only: true},
        { href: '#/moderation', label: 'Модерация', authorized_only: true, guest_only:false},
        { href: '#/logout', label: 'Выйти', authorized_only: true, guest_only:false}
    ];

    for (const navLink of navLinks) {
        if (navLink.authorized_only && !authorized) {
            continue;
        }
        if (navLink.guest_only && authorized){
            continue
        }
        const item = document.createElement('li');
        const link = document.createElement('a');
        link.href = navLink.href;
        link.textContent = navLink.label;
        item.appendChild(link);
        list.appendChild(item);
    }

    nav.appendChild(title);
    nav.appendChild(list);
    header.appendChild(nav);

    return header;
}
