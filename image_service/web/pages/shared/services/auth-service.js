import { apiRequest } from './api-client.js';

export class SessionExpiredError extends Error {
    constructor() {
        super('Сессия истекла');
        this.name = 'SessionExpiredError';
    }
}

let sessionExpiredHandler = null;
let sessionExpiredNotified = false;

export function setSessionExpiredHandler(handler) {
    sessionExpiredHandler = handler;
}

export function resetSessionExpiredState() {
    sessionExpiredNotified = false;
}

export function handleSessionExpired() {
    if (sessionExpiredNotified) {
        return;
    }

    sessionExpiredNotified = true;
    if (typeof sessionExpiredHandler === 'function') {
        sessionExpiredHandler();
    }
}

function isUnauthorizedError(error) {
    return Boolean(error && error.unauthorized);
}

export async function registerUser(data) {
    return apiRequest('/auth/register', {
        method: 'POST',
        body: data,
    });
}

export async function loginUser(data) {
    return apiRequest('/auth/login', {
        method: 'POST',
        body: data,
    });
}

export async function getMe() {
    return apiRequest('/auth/me', {
        method: 'GET',
    });
}

export async function resolveAuthState() {
    try {
        await getMe();
        return { authorized: true, sessionExpired: false };
    } catch {
        return { authorized: false, sessionExpired: true };
    }

}

export async function logoutUser() {
    return apiRequest('/auth/logout', {
        method: 'POST',
    });
}
