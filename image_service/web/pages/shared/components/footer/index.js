const footerStylesHref = new URL('./footer.css', import.meta.url).href;

function ensureFooterStyles() {
    if (document.querySelector('link[data-component="footer"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = footerStylesHref;
    link.dataset.component = 'footer';
    document.head.appendChild(link);
}

export function renderFooter() {
    ensureFooterStyles();

    const footer = document.createElement('footer');
    footer.className = 'site-footer';
    const text = document.createElement('p');
    text.textContent = `© ${new Date().getFullYear()} Image Service`;

    footer.appendChild(text);

    return footer;
}
