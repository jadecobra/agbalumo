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

document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('dialog[id][data-auto-open]').forEach(modal => {
        modal.showModal();
    });

    document.querySelectorAll('dialog[id]').forEach(modal => {
        modal.addEventListener('click', function(event) {
            const rect = modal.getBoundingClientRect();
            const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height
              && rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                modal.close();
            }
        });
    });
});
