const DEFAULT_API_BASE_URL = '/api';

export const API_BASE_URL = DEFAULT_API_BASE_URL.replace(/\/+$/, '');

function createErrorMessage(responseStatus, payload) {
    if (payload && typeof payload.err_message === 'string' && payload.err_message !== '') {
        return payload.err_message;
    }

    return `Request failed with status ${responseStatus}`;
}

export async function apiRequest(path, options = {}) {
    const { method = 'GET', body } = options;
    const response = await fetch(`${API_BASE_URL}${path}`, {
        method,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
        },
        body: body ? JSON.stringify(body) : undefined,
    });

    const isJson = response.headers.get('content-type')?.includes('application/json');
    const payload = isJson ? await response.json() : null;

    if (!response.ok || (payload && payload.success === false)) {
        const error = new Error(createErrorMessage(response.status, payload));
        error.status = response.status;
        error.unauthorized = response.status === 401;
        throw error;
    }

    return payload ? payload.data : null;
}
