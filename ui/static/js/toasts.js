function showToast(message, type, duration) {
    if (!duration) duration = 5000;
    const toastId = 'toast-' + Date.now();
    
    const toastHTML = `
        <div id="${toastId}" 
             class="fixed top-4 right-4 z-50 max-w-sm w-full bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-xl shadow-lg p-4 flex items-start gap-3 animate-in slide-in-from-top-2 fade-in"
             role="alert">
            <span class="material-symbols-outlined text-red-500 text-[20px] mt-0.5">error</span>
            <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-red-800 dark:text-red-200">Error</p>
                <p class="text-sm text-red-600 dark:text-red-300 mt-1">${message}</p>
            </div>
            <button onclick="closeToast('${toastId}')" 
                    class="text-red-400 hover:text-red-600 dark:hover:text-red-200 transition-colors">
                <span class="material-symbols-outlined text-[18px]">close</span>
            </button>
        </div>
    `;
    
    document.body.insertAdjacentHTML('beforeend', toastHTML);
    
    setTimeout(() => {
        dismissToast(toastId);
    }, duration);
}

function closeToast(toastId) {
    dismissToast(toastId);
}

function dismissToast(toastId) {
    const toast = document.getElementById(toastId);
    if (toast) {
        toast.style.animation = 'fade-out 0.3s ease-out forwards';
        setTimeout(() => toast.remove(), 300);
    }
}

function initToasts() {
    document.querySelectorAll('[data-toast-auto-dismiss]').forEach(toast => {
        const id = toast.id;
        const duration = parseInt(toast.dataset.toastAutoDismiss, 10) || 5000;
        setTimeout(() => dismissToast(id), duration);
    });
}

document.addEventListener('DOMContentLoaded', function() {
    window.closeToast = closeToast;
    window.dismissToast = dismissToast;
    initToasts();
});

document.body.addEventListener('htmx:afterSwap', initToasts);
