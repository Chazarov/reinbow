import {loadMainFeed, loadUserFeedById, reloadFeed} from '../feed/feed.js';
import noImage from '../../img/noimage.png';
import {closeModal, showGuestUI} from "../index";
import FeedAPI from "../client/feed_service.js";
/**
 * Переключает видимость изображений в слайдере поста на предыдущее.
 *
 * Находит текущее видимое изображение (с классом `showed-image`) внутри блока
 * с id `images_${postId}` и, если существует предыдущий соседний элемент,
 * скрывает текущее и показывает предыдущее.
 *
 * @function prevImage
 * @param {number|string} postId - Идентификатор поста.
 * @returns {void}
 */
function prevImage(postId) {
    const imagesBlock = document.getElementById("images_" + postId);
    if (!imagesBlock) return;
    const currentImage = imagesBlock.querySelector(".showed-image");
    if (!currentImage) return;
    const prevImageLink = currentImage.parentElement.previousElementSibling;
    if (prevImageLink && prevImageLink.children.length !== 0) {
        const prevImg = prevImageLink.children[0];
        if (prevImg) {
            currentImage.classList.remove('showed-image');
            prevImg.classList.add('showed-image');
        }
    }
}

/**
 * Переключает видимость изображений в слайдере поста на следующее.
 *
 * Находит текущее видимое изображение (с классом `showed-image`) внутри блока
 * с id `images_${postId}` и, если существует следующий соседний элемент,
 * скрывает текущее и показывает следующий.
 *
 * @function nextImage
 * @param {number|string} postId - Идентификатор поста.
 * @returns {void}
 */
function nextImage(postId) {
    const imagesBlock = document.getElementById("images_" + postId);
    if (!imagesBlock) return;
    const currentImage = imagesBlock.querySelector(".showed-image");
    if (!currentImage) return;
    const nextImageLink = currentImage.parentElement.nextElementSibling;
    if (nextImageLink && nextImageLink.children.length !== 0) {
        const nextImg = nextImageLink.children[0];
        if (nextImg) {
            currentImage.classList.remove('showed-image');
            nextImg.classList.add('showed-image');
        }
    }
}

/**
 * Генерирует HTML-разметку карточки поста со слайдером изображений.
 *
 * Использует `process.env.IS_URL` - URL ImageService - для формирования полного URL изображений.
 * Если у поста есть несколько изображений, добавляет кнопки навигации.
 *
 * @function createPost
 * @param {Object} post - Объект с данными поста.
 * @param {number} post.id - Уникальный идентификатор поста.
 * @param {number} [post.user_id] - ID автора поста.
 * @param {string} [post.username] - Имя автора.
 * @param {string} [post.title] - Заголовок поста.
 * @param {string} [post.description] - Описание поста.
 * @param {string} [post.created_at] - Дата создания.
 * @param {Array<number|string>} [post.images] - Массив идентификаторов изображений.
 * @returns {string} HTML-строка с карточкой поста.
 *
 * @example
 * const postHTML = createPost({ id: 1, title: "Закат", images: [101, 102] });
 */
export function createPost(post) {
    let sliderImages = "<div class='img-slider'>";

    for (const imageId of post.images || []) {
        sliderImages += `<a href="${process.env.IS_URL}/api/guest/image/${imageId}" target="_blank">
            <img id="image_${imageId}" src="${process.env.IS_URL}/api/guest/image/${imageId}?type=icon" class="post-image">
        </a>`;
    }
    sliderImages += "</div>";

    const tools = `
        <div class="image-tools">
            <span class="image-btn imagePrev" data-post-id="${post.id}" data-action="prev">
                <span class="fa fa-arrow-left"></span>
            </span>
            <span class="image-btn imageNext" data-post-id="${post.id}" data-action="next">
                <span class="fa fa-arrow-right"></span>
            </span>
        </div>`;

    return `
        <div class="post-card" id="post_${post.id}">
            <div class="imagesSlider">
                <div class="images" id="images_${post.id}">
                    ${sliderImages}
                    ${post.images && post.images.length > 1 ? tools : ''}
                </div>
            </div>
            <div class="post-body">
                <h3 class="post-title">${post.title || 'Без названия'}</h3>
                <p class="post-description">${post.description || ''}</p>
                <div class="post-meta">
                    <span class="post-author" data-user-id="${post.user_id}">${post.username || 'Аноним'}</span>
                    <span class="post-date">${post.created_at}</span>
                </div>
            </div>
        </div>
    `;
}

/**
 * Навешивает обработчики событий на элементы внутри постов после рендеринга.
 *
 * - Первому изображению в каждом слайдере добавляет класс `showed-image`.
 * - Для кнопок `.imagePrev` и `.imageNext` делегирует вызовы `prevImage` и `nextImage`.
 * - Для элементов `.post-author` добавляет клик-обработчик, загружающий ленту автора,
 *   если текущий пользователь подтверждён.
 *
 * @function createPostsHandlers
 * @returns {void}
 */
export function createPostsHandlers() {
    // Показать первые изображения
    document.querySelectorAll('.img-slider a:first-child .post-image').forEach(el => {
        el.classList.add("showed-image");
    });

    // Ловим все ошибки загрузки картинок на странице на фазе перехвата
    document.addEventListener('error', (event) => {
        if (event.target.tagName.toLowerCase() === 'img') {
            event.target.src = noImage; // Подставляем заглушку
        }
    }, true); // <- Важно: true включает перехват (capture)


    // Делегирование навигации слайдера
    document.querySelectorAll('.imagePrev, .imageNext').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const postId = e.currentTarget.dataset.postId;
            const action = e.currentTarget.dataset.action;
            if (action === 'prev') {
                prevImage(postId);
            } else if (action === 'next') {
                nextImage(postId);
            }
        });
    });

    // Обработчики клика по автору
    document.querySelectorAll('.post-author').forEach(el => {
        el.addEventListener('click', async (e) => {
            const user = await FeedAPI.getUserData();
            if (!user) {
                showGuestUI();
                return;
            }
            const userId = e.target.dataset.userId;
            if (userId && user.data.is_confirmed) {
                loadUserFeedById(parseInt(userId));
            }
        });
    });
}