'use strict';
// Server-rendered forms also work without JavaScript. This only narrows choices
// and displays progress; validation and duplicate protection live on the server.
for (const form of document.querySelectorAll('[data-kind-filter]')) {
  const kind = form.elements.namedItem('kind');
  const category = form.elements.namedItem('category');
  const update = () => {
    for (const option of category.options) {
      const mismatch = option.dataset.kind && option.dataset.kind !== kind.value;
      option.hidden = Boolean(mismatch);
      option.disabled = Boolean(mismatch);
      if (mismatch && option.selected) category.value = '';
    }
  };
  kind.addEventListener('change', update);
  update();
}
for (const form of document.querySelectorAll('form[method="post"]')) {
  form.addEventListener('submit', () => {
    const message = document.createElement('p');
    message.className = 'busy-note';
    message.setAttribute('role', 'status');
    message.textContent = 'Сохраняем…';
    form.append(message);
  }, {once: true});
}
