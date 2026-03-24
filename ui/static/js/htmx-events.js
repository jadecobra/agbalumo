// HTMX Component Integration & CSRF Token handling
document.addEventListener('DOMContentLoaded', () => {
    setupHtmxIntegration();
    setupCsrf();
});

function setupHtmxIntegration() {
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        // htmx:afterSwap provides the swapped element in evt.detail.elt
        const elt = evt.detail.elt;
        if (!elt) return;

        // Auto-open and init edit modals inserted via HTMX
        function initEditDialog(d) {
            if (d.dataset.autoOpen === 'true' && !d.open) {
                d.showModal();
                const listingId = d.dataset.listingId;
                if (listingId) {
                    if (window.initEditTypeToggle) window.initEditTypeToggle(listingId);
                    if (window.initEditMaps) window.initEditMaps(listingId);
                    if (window.initEditImagePreview) window.initEditImagePreview(listingId);
                }
            }
        }

        // Check if the swapped element is the dialog itself
        if (elt.tagName === 'DIALOG' && elt.id.startsWith('edit-listing-modal-')) {
            initEditDialog(elt);
        } else {
            // Or if the dialog is contained within the swapped element
            const dialogs = elt.querySelectorAll ? elt.querySelectorAll('dialog[id^="edit-listing-modal-"]') : [];
            dialogs.forEach(d => initEditDialog(d));
        }

        // Init create image preview when create modal is swapped in
        if (elt.id === 'create-listing-modal' || (elt.querySelectorAll && elt.querySelector('#create-listing-modal'))) {
            if (window.initCreateImagePreview) window.initCreateImagePreview();
        }
    });

    // Auto-close edit modal after successful save
    document.body.addEventListener('htmx:afterRequest', (evt) => {
        const form = evt.detail.elt;
        if (!form || !form.id || !form.id.startsWith('edit-form-')) return;
        if (!evt.detail.successful) return;

        const dialog = form.closest('dialog');
        if (dialog) {
            dialog.close();
            dialog.remove();
        }
    });

    // Initialize custom dropdowns on HTMX swap
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        if (typeof initCustomDropdownsActiveState === 'function') {
            initCustomDropdownsActiveState(evt.detail.elt);
        }
    });
}

function setupCsrf() {
    document.body.addEventListener('htmx:configRequest', (evt) => {
        const token = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
        if (token) {
            evt.detail.headers['X-CSRF-Token'] = token;
        }
    });
}
