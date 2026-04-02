function setupAuthActions() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-auth-action="modal"]');
        if (btn) {
            const isAuthenticated = btn.dataset.userAuthenticated === 'true';
            if (isAuthenticated) {
                const modalId = btn.dataset.modalId;
                const modal = document.getElementById(modalId);
                if (modal) modal.showModal();
            } else {
                window.location.href = '/auth/google/login';
            }
        }
    });
}
