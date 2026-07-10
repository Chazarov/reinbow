import { renderHeader } from './pages/shared/components/header/index.js';
import { renderFooter } from './pages/shared/components/footer/index.js';
import { renderRegisterPage } from './pages/register/index.js';
import { renderLoginPage } from './pages/login/index.js';
import { renderModerationPage } from './pages/moderation/index.js';
import {
    handleSessionExpired,
    logoutUser,
    resetSessionExpiredState,
    resolveAuthState,
    setSessionExpiredHandler,
} from './pages/shared/services/auth-service.js';

const app = document.getElementById('app');
const main = document.createElement('main');
main.id = 'main';

const routes = {
    '/register': () =>
        renderRegisterPage({
            onRegisterSuccess: () => navigate('/login'),
        }),
    '/login': () =>
        renderLoginPage({
            onLoginSuccess: () => {
                resetSessionExpiredState();
                navigate('/moderation');
            },
        }),
    '/moderation': () => renderModerationPage(),
};

const protectedRoutes = new Set(['/moderation']);

setSessionExpiredHandler(() => {
    alert('Сессия истекла');
    navigate('/login');
});

function navigate(route) {
    window.location.hash = route;
}

function getCurrentRoute() {
    const hash = window.location.hash.slice(1);
    return hash || '/login';
}

async function renderRoute() {
    const route = getCurrentRoute();
    const isLogoutRoute = route === '/logout';
    const isProtectedRoute = protectedRoutes.has(route);
    const { authorized, sessionExpired } = isLogoutRoute
        ? { authorized: false, sessionExpired: false }
        : await resolveAuthState();

    const oldHeader = document.querySelector('.site-header');
    if (oldHeader) {
        oldHeader.replaceWith(renderHeader({ authorized }));
    }

    if (isLogoutRoute) {
        try {
            await logoutUser();
        } catch (error) {
            console.error('Logout failed:', error);
        }

        navigate('/login');
        return;
    }

    if (isProtectedRoute && !authorized) {
        if (sessionExpired) {
            handleSessionExpired();
        } else {
            navigate('/login');
        }
        return;
    }

    const pageFactory = routes[route] || routes['/login'];
    main.replaceChildren(pageFactory());
}

app.appendChild(renderHeader());
app.appendChild(main);
app.appendChild(renderFooter());

window.addEventListener('hashchange', renderRoute);
renderRoute();
