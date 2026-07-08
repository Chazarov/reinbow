export function fetchWithAuth(url, options = {}) {
    const token = localStorage.getItem('token');
    if (token) {
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`
        };
    }
    return fetch(url, options);
}

export async function login(loginData) {
    const response = await fetch("/api/moderator_login", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(loginData)
    });
    return response.json();
}

export async function getNotifySettings() {
    const response = await fetchWithAuth("/api/get_notify_settings");
    return response.json();
}

export async function saveNotifySettings(payload) {
    return fetchWithAuth("/api/notify_types", {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    });
}