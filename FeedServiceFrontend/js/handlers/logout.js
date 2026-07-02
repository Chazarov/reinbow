export function clearSession() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
}

export function initLogoutHandler() {
    const btnLogout = document.getElementById('btnLogout');
    
    btnLogout.addEventListener('click', async () => {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
            try {
                await fetch(process.env.FS_URL + '/logout', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken })
                });
            } catch (e) {
            }
        }
        clearSession();
        showGuestUI();
    });
}