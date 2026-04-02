function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.showModal();
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.close();
        modal.remove();
    }
}

document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('dialog[id][data-auto-open]').forEach(modal => {
        modal.showModal();
    });

    document.querySelectorAll('dialog[id]').forEach(modal => {
        modal.addEventListener('click', function (event) {
            const rect = modal.getBoundingClientRect();
            const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height
                && rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                modal.close();
            }
        });
    });

    // Auto-open detail modals
    document.querySelectorAll('dialog[id^="detail-modal-"]').forEach(modal => {
        modal.showModal();
    });

    // Auto-open feedback modal
    const feedbackModal = document.getElementById('feedback-modal');
    if (feedbackModal) {
        feedbackModal.showModal();
    }

    // Initialize edit modals if present on initial load
    document.querySelectorAll('dialog[id^="edit-listing-modal-"]').forEach(modal => {
        const listingId = modal.id.replace('edit-listing-modal-', '');
        if (listingId) {
            if (typeof initEditTypeToggle === 'function') initEditTypeToggle(listingId);
            if (typeof initEditImagePreview === 'function') initEditImagePreview(listingId);
        }
    });
});

// HTMX Listener for dynamically loaded content
document.addEventListener('htmx:afterSettle', function (event) {
    const target = event.detail.target;
    if (!target) return;

    // Handle new edit listing modals
    const editModal = target.querySelector('dialog[id^="edit-listing-modal-"]') ||
        (target.id?.startsWith('edit-listing-modal-') ? target : null);

    if (editModal) {
        const listingId = editModal.id.replace('edit-listing-modal-', '');
        if (listingId) {
            editModal.showModal();
            if (typeof initEditTypeToggle === 'function') initEditTypeToggle(listingId);
            if (typeof initEditImagePreview === 'function') initEditImagePreview(listingId);
        }
    }

    // Auto-open detail modals loaded via HTMX
    target.querySelectorAll('dialog[id^="detail-modal-"]').forEach(modal => {
        modal.showModal();
    });

    // Auto-open feedback modal loaded via HTMX
    const feedbackModal = target.querySelector('#feedback-modal') ||
        (target.id === 'feedback-modal' ? target : null);
    if (feedbackModal) {
        feedbackModal.showModal();
    }

    // Auto-open any loaded data-auto-open modals
    const autoOpenModals = target.querySelectorAll('dialog[id][data-auto-open]');
    if (autoOpenModals) {
        autoOpenModals.forEach(modal => {
            modal.showModal();
        });
    }
});
