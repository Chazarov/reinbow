const formInputStylesHref = new URL('./form-input.css', import.meta.url).href;

function ensureFormInputStyles() {
    if (document.querySelector('link[data-component="form-input"]')) {
        return;
    }

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = formInputStylesHref;
    link.dataset.component = 'form-input';
    document.head.appendChild(link);
}

export function createFormInput(config) {
    ensureFormInputStyles();

    const { label, name, type = 'text', placeholder = '' } = config;
    const wrapper = document.createElement('label');
    wrapper.className = 'form-input';

    const title = document.createElement('span');
    title.className = 'form-input-label';
    title.textContent = label;

    const input = document.createElement('input');
    input.className = 'form-input-control';
    input.type = type;
    input.name = name;
    input.placeholder = placeholder;
    input.required = true;

    wrapper.appendChild(title);
    wrapper.appendChild(input);

    return { wrapper, input };
}
