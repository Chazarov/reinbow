import {clearSession} from "./logout";

export async function refreshAccessToken() {
    try {
        const response = await fetch('/api/refresh', {
            method: 'GET',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include'
        });

        const data = await response.json();

        if (!response.ok || !data.success) {
            clearSession()
            throw new Error(data.err_message || 'Токен некорректный');
        }

        // data.data должен содержать { access_token, user }
        const {access_token, user} = data.data;
        if (!access_token || !user) {
            clearSession()
            throw new Error('Некорректный ответ сервера');
        }

        // Сохраняем сессию
        saveSession(user);
    } catch (err) {
        clearSession()
        console.log(err.message);
    }
}