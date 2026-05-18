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
                navigator.geolocation.getCurrentPosition(
                    (position) => {
                        if (!silent) hidePrompt();
                        localStorage.setItem(STORAGE_KEY, 'true');
                        
                        const lat = position.coords.latitude;
                        const lng = position.coords.longitude;
                        
                        if (window.htmx) {
                            const urlParams = new URLSearchParams(window.location.search);
                            const values = { 
                                lat: lat, 
                                lng: lng, 
                                radius: 10,
                                start_ts: startTime 
                            };
                            // Preserve existing query params like limit, type, and search queries
                            for (const [key, val] of urlParams.entries()) {
                                if (!(key in values)) {
                                    values[key] = val;
                                }
                            }
                            htmx.ajax('GET', '/listings/fragment', {
                                values: values,
                                target: '#listings-container',
                                swap: 'innerHTML'
                            });
                        }
                    },
                    (error) => {
                        if (!silent) {
                            console.warn("Location access denied or failed:", error);
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
