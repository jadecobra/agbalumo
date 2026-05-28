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

                // Clean up URL query parameters in address bar
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

            // If the HTML permission prompt was dismissed or native permission was previously denied,
            // show the custom HTML permission prompt again.
            const prompt = document.getElementById('location-permission-prompt');
            const isCurrentlyDenied = text && text.textContent && text.textContent.includes('Denied');
            if (prompt && (sessionStorage.getItem('agbalumo_geo_dismissed') === 'true' || isCurrentlyDenied)) {
                sessionStorage.removeItem('agbalumo_geo_dismissed');
                prompt.classList.remove('hidden', 'animate-out', 'fade-out');
                applyDefaultState();
                return;
            }

            // Clear any legacy dismissed state before a fresh gesture attempt.
            // This gives the browser the cleanest chance to show the native prompt on re-click after denial.
            sessionStorage.removeItem('agbalumo_geo_dismissed');

            // Show loading state
            clickBtn.disabled = true;
            if (icon) icon.classList.add('hidden');
            if (spinner) spinner.classList.remove('hidden');
            if (text) text.textContent = 'Locating...';

            if ("geolocation" in navigator) {
                navigator.geolocation.getCurrentPosition(
                    function(position) {
                        const lat = position.coords.latitude;
                        const lng = position.coords.longitude;

                        sessionStorage.setItem('agbalumo_lat', lat);
                        sessionStorage.setItem('agbalumo_lng', lng);

                        applyActiveState();

                        if (window.htmx) {
                            const urlParams = new URLSearchParams(window.location.search);
                            const values = {
                                lat: lat,
                                lng: lng,
                                radius: window.filterState?.radius || '10'
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
                    function(error) {
                        const isDenied = error.code === 1 || error.code === error.PERMISSION_DENIED;
                        let errMsg = 'Error: Unavailable';
                        if (isDenied) {
                            errMsg = 'Denied - tap to retry';
                        }
                        if (text) text.textContent = errMsg;
                        if (spinner) spinner.classList.add('hidden');

                        // Re-enable immediately so a follow-up click is a clean user gesture
                        // (gives the browser the best chance to re-show the native prompt).
                        clickBtn.disabled = false;

                        if (!isDenied) {
                            // Transient errors (timeout/unavailable) can flash and reset.
                            setTimeout(function() {
                                applyDefaultState();
                            }, 2000);
                        }
                        // For denial: leave the button enabled in a retryable state.
                        // Next click will attempt getCurrentPosition again (fresh gesture).
                    },
                    { enableHighAccuracy: false, timeout: 15000, maximumAge: 30000 }
                );
            } else {
                if (text) text.textContent = 'Error: Unsupported';
                if (spinner) spinner.classList.add('hidden');
                setTimeout(function() {
                    applyDefaultState();
                }, 2000);
            }
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
