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

        // Check for existing permission
        if (navigator.permissions && navigator.permissions.query) {
            navigator.permissions.query({ name: 'geolocation' }).then(result => {
                if (result.state === 'granted') {
                    // Auto-trigger if already granted
                    triggerGeolocation(true);
                } else if (result.state === 'prompt') {
                    if (!localStorage.getItem(STORAGE_KEY)) {
                        setTimeout(() => {
                            prompt.classList.remove('hidden');
                        }, 1500);
                    }
                }
            });
        } else if (!localStorage.getItem(STORAGE_KEY)) {
            // Fallback for browsers that don't support permissions.query
            setTimeout(() => {
                prompt.classList.remove('hidden');
            }, 1500);
        }

        function triggerGeolocation(silent = false) {
            if ("geolocation" in navigator) {
                const startTime = Date.now();
                if (!silent) {
                    allowBtn.disabled = true;
                    allowBtn.innerText = "FINDING SPOTS...";
                    allowBtn.classList.add('opacity-75', 'cursor-not-allowed');
                }
                navigator.geolocation.getCurrentPosition(
                    (position) => {
                        if (!silent) {
                            allowBtn.innerText = "FOUND!";
                            setTimeout(hidePrompt, 300);
                        } else {
                            hidePrompt();
                        }
                        localStorage.setItem(STORAGE_KEY, 'true');
                        
                        const lat = position.coords.latitude;
                        const lng = position.coords.longitude;
                        
                        if (window.htmx) {
                            htmx.ajax('GET', '/listings/fragment', {
                                values: { 
                                    lat: lat, 
                                    lng: lng, 
                                    radius: 10,
                                    start_ts: startTime 
                                },
                                target: '#listings-container',
                                swap: 'innerHTML',
                                indicator: '#listings-loading'
                            });
                        }
                    },
                    (error) => {
                        if (!silent) {
                            console.warn("Location access denied or failed:", error);
                            allowBtn.innerText = "ALLOW LOCATION";
                            allowBtn.disabled = false;
                            allowBtn.classList.remove('opacity-75', 'cursor-not-allowed');
                            hidePrompt();
                        }
                        localStorage.setItem(STORAGE_KEY, 'true');
                    }
                );
            }
        }

        allowBtn.addEventListener('click', () => triggerGeolocation(false));

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
