function setupAuthActions() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-auth-action="modal"]');
        if (btn) {
            console.log('Auth action triggered:', btn.dataset);
            e.preventDefault();
            e.stopPropagation();
            const isAuthenticated = btn.dataset.userAuthenticated === 'true';
            if (isAuthenticated) {
                const modalId = btn.dataset.modalId;
                const modal = document.getElementById(modalId);
                if (modal) modal.showModal();
            } else {
                const loginModal = document.getElementById('login-prompt-modal');
                if (loginModal) loginModal.showModal();
            }
        }
    });
}
