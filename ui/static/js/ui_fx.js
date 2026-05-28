function shouldRemoveModal(id) {
    if (!id) return false;
    return id.startsWith('detail-modal-') || 
           id.startsWith('edit-listing-modal-') || 
           id.startsWith('edit-modal-') || 
           id === 'feedback-modal' || 
           id === 'profile-modal';
}

function setupModalClosing() {
    document.addEventListener('click', (event) => {
        if (event.target.tagName === 'DIALOG') {
            const dialog = event.target;
            const rect = dialog.getBoundingClientRect();
            const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height &&
                rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                dialog.close();
                if (shouldRemoveModal(dialog.id)) {
                    dialog.remove();
                }
            }
        }
    });
}

function setupModalActions() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-action="close"]');
        if (btn) {
            const modal = btn.closest('dialog') || document.getElementById(btn.dataset.modalId);
            if (modal) {
                modal.close();
                const form = modal.querySelector('form');
                if (form) form.reset();
                if (shouldRemoveModal(modal.id)) {
                    modal.remove();
                }
            }
        }
    });

    document.addEventListener('submit', (e) => {
        const form = e.target;
        if (form.hasAttribute('data-modal-action-submit-close')) {
            const modal = form.closest('dialog');
            if (modal) {
                modal.addEventListener('htmx:afterRequest', () => {
                    modal.close();
                    form.reset();
                }, { once: true });
            }
        }
    });
}

function setupDynamicBackgrounds() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const el = entry.target;
                const url = el.dataset.bgImage;
                if (url) {
                    el.style.backgroundImage = `url('${url}')`;
                    el.style.backgroundSize = 'cover';
                    el.style.backgroundPosition = 'center';
                }
                observer.unobserve(el);
            }
        });
    });

    document.querySelectorAll('.dynamic-bg').forEach(el => observer.observe(el));

    document.body.addEventListener('htmx:afterSwap', (evt) => {
        evt.detail.elt.querySelectorAll('.dynamic-bg').forEach(el => observer.observe(el));
    });
}
