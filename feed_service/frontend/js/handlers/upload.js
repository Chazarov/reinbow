/**
 * Инициализирует обработчики событий для модального окна загрузки постов с изображениями.
 *
 * Привязывает события:
 * - Открытие модального окна по клику на кнопку `#btnUpload`.
 * - Закрытие по клику на элемент `.close` внутри модального окна.
 * - Отображение списка выбранных файлов с валидацией при изменении `#fileInput`.
 * - Асинхронная отправка формы `#uploadForm` с проверкой полей.
 *
 * Правила валидации файлов:
 * - Не более 3 изображений.
 * - Размер каждого файла ≤ 2 МБ.
 * - Допускаются только MIME-типы, начинающиеся с `image/`.
 *
 * При успешной загрузке вызывает:
 * - {@link module:main.closeModal} для закрытия модального окна.
 * - {@link module:feed/feed.reloadFeed} для обновления ленты.
 *
 * Ошибки выводятся в элемент `#uploadError`.
 *
 * @function initUploadHandlers
 * @requires module:main.openModal
 * @requires module:main.closeModal
 * @requires module:feed/feed.reloadFeed
 * @requires HTML-элементы с id: `btnUpload`, `uploadModal`, `uploadForm`,
 *           `fileInput`, `fileList`, `uploadError`, `postTitle`, `postDescription`.
 * @returns {void}
 *
 * @example
 * // Вызов после загрузки DOM
 * document.addEventListener('DOMContentLoaded', initUploadHandlers);
 */
import {openModal, closeModal, showGuestUI} from '../index.js';
import {loadMainFeed, reloadFeed} from '../feed/feed.js';
import FeedAPI from '../client/feed_service.js';
import { ApiError } from '../client/feed_service.js';

export function initUploadHandlers() {
    const uploadBtn = document.getElementById('btn-upload');
    const uploadModal = document.getElementById('uploadModal');
    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const fileList = document.getElementById('fileList');
    const uploadError = document.getElementById('uploadError');

    uploadBtn.addEventListener('click', () => {
        openModal(uploadModal);
    });

    const closeBtn = uploadModal.querySelector('.close');
    if (closeBtn) {
        closeBtn.addEventListener('click', () => closeModal(uploadModal));
    }

    // Отображение выбранных файлов
    fileInput.addEventListener('change', () => {
        fileList.innerHTML = '';
        uploadError.textContent = '';

        const files = Array.from(fileInput.files);
        if (files.length > 3) {
            uploadError.textContent = 'Можно выбрать не более 3 изображений';
            fileInput.value = '';
            return;
        }

        const maxSize = 2 * 1024 * 1024;
        for (const file of files) {
            if (file.size > maxSize) {
                uploadError.textContent = `Файл "${file.name}" превышает 2 МБ`;
                fileInput.value = '';
                fileList.innerHTML = '';
                return;
            }
            if (!file.type.startsWith('image/')) {
                uploadError.textContent = `Файл "${file.name}" не является изображением`;
                fileInput.value = '';
                fileList.innerHTML = '';
                return;
            }
            const item = document.createElement('div');
            item.className = 'file-item';
            item.textContent = `${file.name} (${(file.size / 1024).toFixed(1)} КБ)`;
            fileList.appendChild(item);
        }
    });

    // Отправка
    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        uploadError.textContent = '';
        e.submitter.disabled = true;

        const title = document.getElementById('postTitle').value.trim();
        const description = document.getElementById('postDescription').value.trim();
        const files = fileInput.files;

        if (!title) {
            uploadError.textContent = 'Введите заголовок';
            return;
        }
        if (files.length === 0) {
            uploadError.textContent = 'Выберите хотя бы одно изображение';
            return;
        }

        const formData = new FormData();
        formData.append('title', title);
        formData.append('description', description);
        for (const file of files) {
            formData.append('images', file);
        }

        try {
            await FeedAPI.upload(formData);
            closeModal(uploadModal);
            if (reloadFeed) reloadFeed();
        } catch (err) {
            if (err instanceof ApiError && err.status === 401) {
                closeModal(document.getElementById("uploadModal"));
                loadMainFeed();
                showGuestUI();
            } else {
                uploadError.textContent = err.message || 'Неизвестная ошибка загрузки';
                e.submitter.disabled = false;
            }
        }
        e.submitter.disabled = false;
    });
}