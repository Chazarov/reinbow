import './utils/utils.css';
import './notifications/notifications.js';
import './auth/auth.js';
import './auth/auth.css';    // <-- добавить
import './settings/settings.js';
import './webhooks/webhooks.js';

import { restoreAuthState, initAuthListeners } from './auth/auth.js';
import { initSettingsListeners } from './settings/settings.js';

restoreAuthState();

if (!window._listenersAdded) {
    window._listenersAdded = true;

    document.addEventListener('DOMContentLoaded', function() {
        initAuthListeners();
        initSettingsListeners();
    });
}