(function() {
    const STORAGE_KEY = 'agbalumo_geo_dismissed';
    const PROMPT_ID = 'location-permission-prompt';
    const ALLOW_BTN_ID = 'location-allow-btn';
    const DISMISS_BTN_ID = 'location-dismiss-btn';

    function init() {
        const prompt = document.getElementById(PROMPT_ID);
        const allowBtn = document.getElementById(ALLOW_BTN_ID);
        const dismissBtn = document.getElementById(DISMISS_BTN_ID);

        if (!prompt || !allowBtn || !dismissBtn) return;

        // Check if already dismissed or if permission already granted/denied
        if (localStorage.getItem(STORAGE_KEY)) return;

        // Show prompt with a small delay for better UX
        setTimeout(() => {
            prompt.classList.remove('hidden');
        }, 1500);

        allowBtn.addEventListener('click', () => {
            if ("geolocation" in navigator) {
                navigator.geolocation.getCurrentPosition(
                    (position) => {
                        console.log("Location access granted:", position);
                        // Hide prompt and mark as dismissed so it doesn't show again
                        hidePrompt();
                        localStorage.setItem(STORAGE_KEY, 'true');
                        
                        // Future: trigger HTMX search with lat/lon
                        // For now, we just follow the "clean UX" requirement
                    },
                    (error) => {
                        console.warn("Location access denied or failed:", error);
                        hidePrompt();
                        // Mark as dismissed to avoid intrusion
                        localStorage.setItem(STORAGE_KEY, 'true');
                    }
                );
            } else {
                hidePrompt();
                localStorage.setItem(STORAGE_KEY, 'true');
            }
        });

        dismissBtn.addEventListener('click', () => {
            hidePrompt();
            localStorage.setItem(STORAGE_KEY, 'true');
        });
    }

    function hidePrompt() {
        const prompt = document.getElementById(PROMPT_ID);
        if (prompt) {
            prompt.classList.add('animate-out', 'fade-out');
            // Wait for animation to finish before adding hidden
            setTimeout(() => {
                prompt.classList.add('hidden');
            }, 500);
        }
    }

    // Initialize on DOM load
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
