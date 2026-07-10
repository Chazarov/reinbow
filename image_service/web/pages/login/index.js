import { createFormInput } from '../shared/components/form-input/index.js';
import { loginUser } from '../shared/services/auth-service.js';

const loginStylesHref = new URL('./login.css', import.meta.url).href;

function ensureLoginStyles() {
    if (document.querySelector('link[data-page="login"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = loginStylesHref;
    link.dataset.page = 'login';
    document.head.appendChild(link);
}

export function renderLoginPage(config = {}) {
    ensureLoginStyles();

    const { onLoginSuccess } = config;
    const container = document.createElement('section');
    container.className = 'auth-page';

    const card = document.createElement('article');
    card.className = 'auth-card';

    const title = document.createElement('h1');
    title.textContent = 'Вход';

    const form = document.createElement('form');
    form.className = 'auth-form';

    const usernameField = createFormInput({
        label: 'Логин',
        name: 'login',
        placeholder: 'Введите логин',
    });

    const passwordField = createFormInput({
        label: 'Пароль',
        name: 'pass',
        type: 'password',
        placeholder: 'Введите пароль',
    });

    const submit = document.createElement('button');
    submit.type = 'submit';
    submit.className = 'auth-submit';
    submit.textContent = 'Войти';

    const message = document.createElement('p');
    message.className = 'auth-message';

    form.appendChild(usernameField.wrapper);
    form.appendChild(passwordField.wrapper);
    form.appendChild(submit);
    form.appendChild(message);

    form.addEventListener('submit', async (event) => {
        event.preventDefault();
        message.classList.remove('is-success');
        message.textContent = '';
        submit.disabled = true;

        try {
            await loginUser({
                login: usernameField.input.value,
                pass: passwordField.input.value,
            });

            message.classList.add('is-success');
            message.textContent = 'Успешный вход.';
            if (typeof onLoginSuccess === 'function') {
                onLoginSuccess();
            }
        } catch (error) {
            message.classList.remove('is-success');
            message.textContent = error.message;
        } finally {
            submit.disabled = false;
        }
    });

    card.appendChild(title);
    card.appendChild(form);
    container.appendChild(card);

    return container;
}
