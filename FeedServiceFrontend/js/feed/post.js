/**
 * Переключает видимость изображений в слайдере поста на предыдущее.
 *
 * Находит текущее видимое изображение (с классом `showedImage`) внутри блока
 * с id `images_${postId}` и, если существует предыдущий соседний элемент,
 * скрывает текущее и показывает предыдущее.
 *
 * @global
 * @function prevImage
 * @memberof window
 * @param {number|string} postId - Идентификатор поста, для которого переключается слайдер.
 * @requires HTML-элемент с id `images_${postId}`, внутри которого есть элементы
 *           с классом `showedImage`.
 * @returns {void}
 *
 * @example
 * // Вызов из обработчика клика
 * <span onClick="prevImage(42)">←</span>
 */
window.prevImage = (postId) => {
    const imagesBlock = document.getElementById("images_" + postId);
    const currentImage = imagesBlock.querySelector(".showedImage");
    const prevImage = currentImage.previousElementSibling;
    if(prevImage) {
        currentImage.classList.remove('showedImage');
        prevImage.classList.add('showedImage');
    }
}

/**
 * Переключает видимость изображений в слайдере поста на следующее.
 *
 * Находит текущее видимое изображение (с классом `showedImage`) внутри блока
 * с id `images_${postId}` и, если существует следующий соседний элемент,
 * скрывает текущее и показывает следующий.
 *
 * @global
 * @function nextImage
 * @memberof window
 * @param {number|string} postId - Идентификатор поста, для которого переключается слайдер.
 * @requires HTML-элемент с id `images_${postId}`, внутри которого есть элементы
 *           с классом `showedImage`.
 * @returns {void}
 *
 * @example
 * // Вызов из обработчика клика
 * <span onClick="nextImage(42)">→</span>
 */
window.nextImage = (postId) => {
    const imagesBlock = document.getElementById("images_" + postId);
    const currentImage = imagesBlock.querySelector(".showedImage");
    const nextImage = currentImage.nextElementSibling;
    if(nextImage) {
        currentImage.classList.remove('showedImage');
        nextImage.classList.add('showedImage');
    }
}

/**
 * Генерирует HTML-разметку карточки поста со слайдером изображений.
 *
 * Использует `process.env.IS_URL` - URL ImageService - для формирования полного URL изображений.
 * Если у поста есть несколько изображений, добавляет кнопки навигации,
 * вызывающие глобальные функции {@link window.prevImage} и {@link window.nextImage}.
 *
 * @function createPost
 * @param {Object} post - Объект с данными поста.
 * @param {number} post.id - Уникальный идентификатор поста.
 * @param {number} [post.user_id] - ID автора поста.
 * @param {string} [post.username] - Имя автора.
 * @param {string} [post.title] - Заголовок поста.
 * @param {string} [post.description] - Описание поста.
 * @param {string} [post.created_at] - Дата создания (в строковом представлении).
 * @param {Array<number|string>} [post.images] - Массив идентификаторов изображений.
 * @returns {string} HTML-строка с карточкой поста.
 * @requires process.env.IS_URL - базовый URL ImageService для загрузки изображений.
 *
 * @example
 * const postHTML = createPost({
 *   id: 1,
 *   title: "Закат",
 *   images: [101, 102],
 *   username: "Анна",
 *   created_at: "2025-01-01"
 * });
 * document.body.innerHTML += postHTML;
 */
export function createPost(post) {
    let sliderImages = "<div>"

    for(const imageId of post.images || []) {
        sliderImages += `<a href="${process.env.IS_URL}/api/images/${imageId}" target="_blank"><img alt="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSPshUBFBR7xshIaXFc_ir-eagtAueBv7aX5Nvxdny6qg&s" 
            id="image_${imageId}"  src="${process.env.IS_URL}/api/images/icon/${imageId}"" class="postImage"></a>`
    }
    sliderImages += "</div>"

    const tools = `
        <div class="imageTools">
            <span onClick="prevImage(${post.id})" class="imageBtn imagePrev">&larr;</span>
            <span onClick="nextImage(${post.id})" class="imageBtn imageNext">&rarr;</span>
        </div>`


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

export function createPostsHandlers() {
    document.querySelectorAll('.postImage:first-child').forEach(el => {
        el.classList.add("showedImage")
    });

    document.querySelectorAll('.post-author').forEach(el => {
        el.addEventListener('click', (e) => {
            const user = getSavedUser()
            const userId = e.target.dataset.userId;
            if (userId) {
                if(user && user.is_confirmed) loadUserFeed(parseInt(userId));
            }
        });
    });
}