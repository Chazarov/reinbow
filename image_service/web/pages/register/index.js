import { createFormInput } from '../shared/components/form-input/index.js';
import { registerUser } from '../shared/services/auth-service.js';

const registerStylesHref = new URL('./register.css', import.meta.url).href;

function ensureRegisterStyles() {
    if (document.querySelector('link[data-page="register"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = registerStylesHref;
    link.dataset.page = 'register';
    document.head.appendChild(link);
}

export function renderRegisterPage(config = {}) {
    ensureRegisterStyles();

    const { onRegisterSuccess } = config;

    const container = document.createElement('section');
    container.className = 'auth-page';

    const card = document.createElement('article');
    card.className = 'auth-card';

    const title = document.createElement('h1');
    title.textContent = 'Регистрация';

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
    submit.textContent = 'Зарегистрироваться';

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
            await registerUser({
                login: usernameField.input.value,
                pass: passwordField.input.value,
            });

            form.reset();
            if (typeof onRegisterSuccess === 'function') {
                onRegisterSuccess();
                return;
            }

            message.classList.add('is-success');
            message.textContent = 'Пользователь создан. Теперь войдите в систему.';
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
