import { API_BASE_URL, apiRequest } from './api-client.js';

export async function getNextUnmoderatedImage() {
    const data = await apiRequest('/image/unmoderated?page=0&page_size=1');
    const images = Array.isArray(data?.images) ? data.images : [];
    return {
        image: images.length > 0 ? images[0] : null,
        total_count: data.total_count ?? 0,
    };
}

export function getImageContentUrl(imageId) {
    return `${API_BASE_URL}/admin/image/${imageId}`;
}

export async function approveImage(imageId) {
    return apiRequest(`/image/${imageId}/approve`, {
        method: 'PUT',
    });
}

export async function banImage(imageId) {
    return apiRequest(`/image/${imageId}/block`, {
        method: 'PUT',
    });
}
