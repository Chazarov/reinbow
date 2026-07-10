import {
    approveImage,
    banImage,
    getImageContentUrl,
    getNextUnmoderatedImage,
} from '../shared/services/image-service.js';
import { SessionExpiredError } from '../shared/services/auth-service.js';

const moderationStylesHref = new URL('./moderation.css', import.meta.url).href;

function ensureModerationStyles() {
    if (document.querySelector('link[data-page="moderation"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = moderationStylesHref;
    link.dataset.page = 'moderation';
    document.head.appendChild(link);
}

export function renderModerationPage() {
    ensureModerationStyles();

    const container = document.createElement('section');
    container.className = 'moderation-page';

    const title = document.createElement('h1');
    title.className = 'moderation-title';
    title.textContent = 'Модерация';

    const imageBox = document.createElement('div');
    imageBox.className = 'moderation-image-box';

    const image = document.createElement('img');
    image.className = 'moderation-image';
    image.alt = 'Фото для модерации';

    const info = document.createElement('p');
    info.className = 'moderation-info';

    const counter_p = document.createElement('p');
    counter_p.className = 'moderation-total-count';

    const actions = document.createElement('div');
    actions.className = 'moderation-actions';

    const approveButton = document.createElement('button');
    approveButton.type = 'button';
    approveButton.className = 'moderation-action moderation-action-approve';
    approveButton.textContent = 'Approved';

    const banButton = document.createElement('button');
    banButton.type = 'button';
    banButton.className = 'moderation-action moderation-action-ban';
    banButton.textContent = 'BAN!';

    actions.appendChild(approveButton);
    actions.appendChild(banButton);
    imageBox.appendChild(image);
    container.appendChild(title);
    container.appendChild(imageBox);
    container.appendChild(info);
    container.appendChild(counter_p);
    container.appendChild(actions);

    let currentImageId = null;

    function setButtonsDisabled(disabled) {
        approveButton.disabled = disabled;
        banButton.disabled = disabled;
    }

    async function loadNextImage() {
        setButtonsDisabled(true);
        info.textContent = 'Загрузка изображения...';

        try {
            const { image: nextImage, total_count } = await getNextUnmoderatedImage();
            if (!nextImage) {
                currentImageId = null;
                image.removeAttribute('src');
                imageBox.classList.add('is-empty');
                image.style.display ='none'
                counter_p.textContent = ` `
                info.textContent = 'Нет изображений со статусом unmoderated.';
                return;
            }

            currentImageId = nextImage.id;
            imageBox.classList.remove('is-empty');
            image.style.display ='block'
            image.src = `${getImageContentUrl(nextImage.id)}?t=${Date.now()}`;
            info.textContent = `ID: ${nextImage.id} | ${nextImage.name}`;
            counter_p.textContent = `Всего ${total_count} изображений требует модерации`
            setButtonsDisabled(false);
        } catch (error) {
            if (error instanceof SessionExpiredError) {
                return;
            }
            currentImageId = null;
            image.removeAttribute('src');
            imageBox.classList.add('is-empty');
            image.style.display ='none'
            info.textContent = error.message;
        }
    }

    async function moderateCurrentImage(action) {
        if (currentImageId === null) {
            return;
        }

        setButtonsDisabled(true);

        try {
            if (action === 'approve') {
                await approveImage(currentImageId);
            } else {
                await banImage(currentImageId);
            }
            await loadNextImage();
        } catch (error) {
            if (error instanceof SessionExpiredError) {
                return;
            }
            info.textContent = error.message;
            setButtonsDisabled(false);
        }
    }

    approveButton.addEventListener('click', () => {
        moderateCurrentImage('approve');
    });

    banButton.addEventListener('click', () => {
        moderateCurrentImage('ban');
    });

    loadNextImage();

    return container;
}
