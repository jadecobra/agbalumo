(function() {
    const DENIED_KEY = 'agbalumo_geo_denied';
    const PROMPT_ID = 'location-permission-prompt';
    const ALLOW_BTN_ID = 'location-allow-btn';
    const DISMISS_BTN_ID = 'location-dismiss-btn';

    function init() {
        const prompt = document.getElementById(PROMPT_ID);
        const allowBtn = document.getElementById(ALLOW_BTN_ID);
        const dismissBtn = document.getElementById(DISMISS_BTN_ID);

        if (!prompt || !allowBtn || !dismissBtn) return;

        // No auto-show or auto-trigger on load. Per spec:
        // - User must click NEAR ME (no spam).
        // - App must explicitly surface its own modal before any geolocation request.
        // - Only explicit ALLOW click ever calls getCurrentPosition.
        // - Denial persists in localStorage so new tabs respect it (no silent geo).

        function triggerGeolocation() {
            if ("geolocation" in navigator) {
                const startTime = Date.now();
                navigator.geolocation.getCurrentPosition(
                    (position) => {
                        hidePrompt();
                        localStorage.removeItem(DENIED_KEY);

                        const lat = position.coords.latitude;
                        const lng = position.coords.longitude;

                        sessionStorage.setItem('agbalumo_lat', lat);
                        sessionStorage.setItem('agbalumo_lng', lng);

                        if (typeof applyActiveState === 'function') {
                            applyActiveState();
                        } else {
                            const nearMeBtn = document.getElementById('near-me-btn');
                            if (nearMeBtn) {
                                const ACTIVE_CLASSES = ['bg-earth-ochre/20', 'text-earth-ochre', 'hover:bg-earth-ochre/30', 'border-earth-ochre/50'];
                                const DEFAULT_CLASSES = ['bg-earth-sand/30', 'text-text-main', 'hover:bg-earth-sand/50', 'border-earth-clay/10'];
                                DEFAULT_CLASSES.forEach(c => nearMeBtn.classList.remove(c));
                                ACTIVE_CLASSES.forEach(c => nearMeBtn.classList.add(c));
                                const textSpan = document.getElementById('near-me-text');
                                if (textSpan) textSpan.textContent = '📍 Nearby';
                                const spinner = document.getElementById('near-me-spinner');
                                if (spinner) spinner.classList.add('hidden');
                                const icon = document.getElementById('near-me-icon');
                                if (icon) icon.classList.remove('hidden');
                            }
                        }

                        if (window.htmx) {
                            const urlParams = new URLSearchParams(window.location.search);
                            const values = {
                                lat: lat,
                                lng: lng,
                                radius: window.filterState?.radius || '10',
                                start_ts: startTime
                            };
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
                        console.warn("Location access denied or failed:", error);
                        hidePrompt();
                        // Explicit ALLOW click is user acceptance.
                        // Do not treat browser failure after ALLOW as our denial.
                        // Surface retryable state on the button instead.
                        if (typeof applyDeniedState === 'function') {
                            applyDeniedState();
                        } else {
                            // Fallback manual DOM update
                            const nearMeBtn = document.getElementById('near-me-btn');
                            if (nearMeBtn) {
                                const DEFAULT_CLASSES = ['bg-earth-sand/30', 'text-text-main', 'hover:bg-earth-sand/50', 'border-earth-clay/10'];
                                const ACTIVE_CLASSES = ['bg-earth-ochre/20', 'text-earth-ochre', 'hover:bg-earth-ochre/30', 'border-earth-ochre/50'];
                                ACTIVE_CLASSES.forEach(c => nearMeBtn.classList.remove(c));
                                DEFAULT_CLASSES.forEach(c => nearMeBtn.classList.add(c));
                                const textSpan = document.getElementById('near-me-text');
                                if (textSpan) textSpan.textContent = 'Denied - tap to retry';
                                const spinner = document.getElementById('near-me-spinner');
                                if (spinner) spinner.classList.add('hidden');
                                const icon = document.getElementById('near-me-icon');
                                if (icon) icon.classList.remove('hidden');
                                nearMeBtn.disabled = false;
                            }
                        }
                    }
                );
            }
        }

        allowBtn.addEventListener('click', () => triggerGeolocation());

        dismissBtn.addEventListener('click', () => {
            hidePrompt();
            localStorage.setItem(DENIED_KEY, 'true');
        });
    }

    function hidePrompt() {
        const prompt = document.getElementById(PROMPT_ID);
        if (prompt) {
            prompt.classList.add('animate-out', 'fade-out');
            setTimeout(() => {
                prompt.classList.add('hidden');
            }, 500);
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
