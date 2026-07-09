import { createPost, createPostsHandlers } from './post.js';
import { openModal, showGuestUI } from '../index.js';
import FeedAPI from '../client/feed_service.js';

// ---------- Состояние модуля ----------
let currentUserFeed = null;        // null = общая лента, иначе ID пользователя
let feedContainer = null;
let feedLoader = null;
let allPostsBtn = null;
let userNameDisplay = null;
let confirmModal = null;

// ---------- Внутренние утилиты ----------

/**
 * Принимает промис, который резолвится массивом постов.
 * Управляет лоадером, очисткой контейнера и отрисовкой/ошибкой.
 */
async function loadAndRender(postsPromise) {
    feedLoader.style.display = 'block';
    feedContainer.innerHTML = '';
    try {
        const posts = await postsPromise;
        renderPosts(posts);
    } catch (err) {
        feedContainer.innerHTML = `<p class="error-message">${err.message}</p>`;
    } finally {
        feedLoader.style.display = 'none';
    }
}

function renderPosts(posts) {
    if (!posts || posts.length === 0) {
        feedContainer.innerHTML = '<p style="text-align:center;padding:2rem;">Пока нет постов</p>';
        return;
    }
    feedContainer.innerHTML = posts.map(post => createPost(post)).join('');
    createPostsHandlers();
}

// ---------- Публичные функции ленты ----------

export function loadMainFeed() {
    currentUserFeed = null;
    allPostsBtn.style.display = 'none';
    if (userNameDisplay) userNameDisplay.style.display = 'inline';
    loadAndRender(FeedAPI.getMainFeed());
}

function loadUserFeed(userId) {
    currentUserFeed = userId;
    allPostsBtn.style.display = 'inline';
    if (userNameDisplay) userNameDisplay.style.display = 'none';
    loadAndRender(FeedAPI.getUserFeed(userId));
}

export function reloadFeed() {
    if (currentUserFeed === null) {
        loadMainFeed();
    } else {
        loadUserFeed(currentUserFeed);
    }
}

export function loadUserFeedById(userId) {
    loadUserFeed(userId);
}

// ---------- Инициализация ----------

export function initFeed() {
    feedContainer = document.getElementById('feedContainer');
    feedLoader = document.getElementById('feedLoader');
    allPostsBtn = document.getElementById('allPostsBtn');
    userNameDisplay = document.getElementById('user-name-display');
    confirmModal = document.getElementById('confirmModal');
    const logo = document.querySelector('.logo');

    allPostsBtn.addEventListener('click', loadMainFeed);

    if (logo) {
        logo.addEventListener('click', (e) => {
            e.preventDefault();
            loadMainFeed();
        });
    }

    if (userNameDisplay) {
        userNameDisplay.addEventListener('click', async () => {
            const user = await FeedAPI.getUserData();
            if (!user) {
                showGuestUI();
                return;
            }
            if (user.id && user.data.is_confirmed) {
                loadUserFeed(user.id);
            } else if (confirmModal) {
                openModal(confirmModal);
            }
        });
    }

    loadMainFeed();
}