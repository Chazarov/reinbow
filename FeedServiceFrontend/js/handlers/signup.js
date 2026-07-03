export function initSignupHandlers() {
    const signupForm = document.getElementById('signupForm');
    const signupError = document.getElementById('signupError');

    signupForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        signupError.textContent = '';

        const username = document.getElementById('signupUsername').value.trim();
        const email = document.getElementById('signupEmail').value.trim();
        const password = document.getElementById('signupPassword').value;
        const confirm = document.getElementById('signupConfirm').value;

        // Валидация
        if (!username || !email || !password || !confirm) {
            signupError.textContent = 'Заполните все обязательные поля';
            return;
        }

        if (password !== confirm) {
            signupError.textContent = 'Пароли не совпадают';
            return;
        }

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            signupError.textContent = 'Введите корректный email';
            return;
        }

        const payload = {
            username,
            password,
            email,
        };

        try {
            const response = await fetch('/api/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.err_message || 'Ошибка регистрации');
            }

            closeModal(signupModal); // форма сбросится
        } catch (err) {
            signupError.textContent = err.message;
        }
    });
}