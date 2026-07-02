import { createPost, createPostsHandlers } from "./post";

export function initFeed() {
    // Контейнер для постов
    const feedContainer = document.getElementById('feedContainer');
    const feedLoader = document.getElementById('feedLoader');
    const allPostsBtn = document.getElementById('allPostsBtn');
    const logo = document.querySelector('.logo');

    // Текущий режим: null = общая лента, иначе ID пользователя
    let currentUserFeed = null;

    // Функция для получения заголовков авторизации
    function getAuthHeaders() {
        const token = localStorage.getItem('access_token');
        return token ? { 'Authorization': 'Bearer ' + token } : {};
    }

    // Загрузка постов с API
    async function fetchPosts(url) {
        console.log(url)
        feedLoader.style.display = 'block';
        feedContainer.innerHTML = '';
        try {
            const response = await fetch(url, { headers: getAuthHeaders() });
            if (!response.ok) throw new Error('Ошибка загрузки');
            const data = await response.json();
            if (data.success) {
                renderPosts(data.data);
            } else {
                feedContainer.innerHTML = `<p class="error-message">` +
                `{data.err_message || 'Не удалось загрузить посты'}</p>`;
            }
        } catch (err) {
            feedContainer.innerHTML = `<p class="error-message">${err.message}</p>`;
        } finally {
            feedLoader.style.display = 'none';
        }
    }

    // Рендер массива постов
    function renderPosts(posts) {
        if (!posts || posts.length === 0) {
            feedContainer.innerHTML = '<p style="text-align:center;padding:2rem;">'
            + 'Пока нет постов</p>';
            return;
        }

        feedContainer.innerHTML = posts.map(post => createPost(post)).join('');

        createPostsHandlers()
    }

    // Загрузка главной ленты
    function loadMainFeed() {
        currentUserFeed = null;
        allPostsBtn.style.display = 'none';
        fetchPosts(process.env.FS_URL + '/loadMainFeed');
    }

    // Загрузка ленты конкретного пользователя
    window.loadUserFeed = (userId) => {
        currentUserFeed = userId;
        allPostsBtn.style.display = 'inline-block';
        fetchPosts(`${process.env.FS_URL}/loadUserFeed/${userId}`);
    }

    // Обработчики
    logo.addEventListener('click', (e) => {
        e.preventDefault();
        loadMainFeed();
    });

    allPostsBtn.addEventListener('click', () => {
        loadMainFeed();
    });

    // Клик по нику в хедере
    const userNameDisplay = document.getElementById('userNameDisplay');
    if (userNameDisplay) {
        userNameDisplay.addEventListener('click', () => {
            const user = JSON.parse(localStorage.getItem('user') || '{}');
            if (user && user.id) {
                loadUserFeed(user.id);
            }
        });
    }

    // Инициализация: загружаем главную ленту
    loadMainFeed();

    // Экспортируем функцию для вызова из других скриптов
    window.reloadFeed = function() {
        if (currentUserFeed === null) {
            loadMainFeed();
        } else {
            loadUserFeed(currentUserFeed);
        }
    };
};