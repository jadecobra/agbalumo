(function() {
    const ACTIVE_CLASSES = ['bg-earth-ochre/20', 'text-earth-ochre', 'hover:bg-earth-ochre/30', 'border-earth-ochre/50'];
    const DEFAULT_CLASSES = ['bg-earth-sand/30', 'text-text-main', 'hover:bg-earth-sand/50', 'border-earth-clay/10'];

    function getElements() {
        return {
            btn: document.getElementById('near-me-btn'),
            icon: document.getElementById('near-me-icon'),
            spinner: document.getElementById('near-me-spinner'),
            text: document.getElementById('near-me-text')
        };
    }

    function applyActiveState() {
        const { btn, text, icon, spinner } = getElements();
        if (!btn) return;

        DEFAULT_CLASSES.forEach(c => btn.classList.remove(c));
        ACTIVE_CLASSES.forEach(c => btn.classList.add(c));
        
        if (text) {
            text.textContent = '📍 Nearby';
        }
        if (icon) icon.classList.remove('hidden');
        if (spinner) spinner.classList.add('hidden');
        btn.disabled = false;
    }

    function applyDefaultState() {
        const { btn, text, icon, spinner } = getElements();
        if (!btn) return;

        ACTIVE_CLASSES.forEach(c => btn.classList.remove(c));
        DEFAULT_CLASSES.forEach(c => btn.classList.add(c));

        if (text) {
            text.textContent = 'Near Me';
        }
        if (icon) icon.classList.remove('hidden');
        if (spinner) spinner.classList.add('hidden');
        btn.disabled = false;
    }

    function applyDeniedState() {
        const { btn, text, icon, spinner } = getElements();
        if (!btn) return;

        ACTIVE_CLASSES.forEach(c => btn.classList.remove(c));
        DEFAULT_CLASSES.forEach(c => btn.classList.add(c));

        if (text) {
            text.textContent = 'Denied - tap to retry';
        }
        if (icon) icon.classList.remove('hidden');
        if (spinner) spinner.classList.add('hidden');
        btn.disabled = false;
    }

    function initNearMe() {
        const { btn } = getElements();
        if (!btn) return;

        // Seed sessionStorage if URL has coordinates (e.g. server-side geo active)
        const urlParams = new URLSearchParams(window.location.search);
        const urlLat = urlParams.get('lat');
        const urlLng = urlParams.get('lng');
        if (urlLat && urlLng) {
            sessionStorage.setItem('agbalumo_lat', urlLat);
            sessionStorage.setItem('agbalumo_lng', urlLng);
        }

        // Persist active state on load
        const lat = sessionStorage.getItem('agbalumo_lat');
        const lng = sessionStorage.getItem('agbalumo_lng');
        if (lat && lng) {
            applyActiveState();
        }

        btn.addEventListener('click', function() {
            const { btn: clickBtn, icon, spinner, text } = getElements();
            if (!clickBtn) return;

            // Check if already active -> Toggle Off
            const lat = sessionStorage.getItem('agbalumo_lat');
            const lng = sessionStorage.getItem('agbalumo_lng');
            if (lat && lng) {
                sessionStorage.removeItem('agbalumo_lat');
                sessionStorage.removeItem('agbalumo_lng');
                applyDefaultState();

                if (history.replaceState) {
                    const url = new URL(window.location.href);
                    url.searchParams.delete('lat');
                    url.searchParams.delete('lng');
                    url.searchParams.delete('radius');
                    history.replaceState(null, '', url.pathname + url.search);
                }

                if (window.htmx) {
                    const urlParams = new URLSearchParams(window.location.search);
                    const values = {};
                    for (const [key, val] of urlParams.entries()) {
                        if (key !== 'lat' && key !== 'lng' && key !== 'radius') {
                            values[key] = val;
                        }
                    }
                    htmx.ajax('GET', '/listings/fragment', {
                        values: values,
                        target: '#listings-container',
                        swap: 'innerHTML'
                    });
                }
                return;
            }

            // Per clarified spec:
            // - User must click NEAR ME to initiate any location request (no load-time spam).
            // - App must explicitly ask via its own modal before any geolocation call.
            // - Only explicit click on ALLOW LOCATION (in the modal) ever triggers the native permission prompt + getCurrentPosition.
            // - Denial (dismiss or native error) lives in localStorage so new tabs / reloads respect it.
            //
            // Therefore: when inactive, NEAR ME click *always* surfaces the custom permission explainer modal.
            // The modal's ALLOW handler (geolocation.js) owns the actual geo call and success side-effects.
            const prompt = document.getElementById('location-permission-prompt');
            if (prompt) {
                // If button is in transient "Denied - tap to retry" state from a prior failed attempt in this session,
                // clear it so the user sees the clean explainer again.
                if (text && text.textContent && text.textContent.includes('Denied')) {
                    applyDefaultState();
                }
                prompt.classList.remove('hidden', 'animate-out', 'fade-out');
                return;
            }

            // Fallback (modal element missing — should never happen as it is rendered server-side).
            // Do nothing rather than silently geolocate.
            return;
        });
    }

    // Auto-inject Coordinates & Reset on City Search
    document.body.addEventListener('htmx:configRequest', function(evt) {
        if (evt.detail.path && evt.detail.path.includes('/listings/fragment')) {
            const cityInputVal = document.getElementById('filter-city')?.value || '';
            const filterStateCity = (window.filterState && window.filterState.city) || '';
            const city = cityInputVal || filterStateCity || evt.detail.parameters['city'] || '';

            if (city && city.trim() !== '') {
                // Clear coords on explicit city search
                sessionStorage.removeItem('agbalumo_lat');
                sessionStorage.removeItem('agbalumo_lng');
                applyDefaultState();
            } else {
                const lat = sessionStorage.getItem('agbalumo_lat');
                const lng = sessionStorage.getItem('agbalumo_lng');
                if (lat && lng) {
                    evt.detail.parameters['lat'] = lat;
                    evt.detail.parameters['lng'] = lng;
                    if (!evt.detail.parameters['radius']) {
                        evt.detail.parameters['radius'] = window.filterState?.radius || '10';
                    }
                }
            }
        }
    });

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initNearMe);
    } else {
        initNearMe();
    }
})();
