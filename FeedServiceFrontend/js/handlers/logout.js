export function clearSession() {
    localStorage.removeItem('user');
}

export function initLogoutHandler() {
    const btnLogout = document.getElementById('btnLogout');
    
    btnLogout.addEventListener('click', async () => {
        try {
            await fetch('/api/logout', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            });
        } catch (e) {
        }
        clearSession();
        showGuestUI();
    });
}