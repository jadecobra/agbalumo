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
            text.textContent = 'Nearby';
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

    function applyDeniedState(message = 'Location blocked by browser') {
        const { btn, text, icon, spinner } = getElements();
        if (!btn) return;

        ACTIVE_CLASSES.forEach(c => btn.classList.remove(c));
        DEFAULT_CLASSES.forEach(c => btn.classList.add(c));

        if (text) {
            text.textContent = message;
        }
        if (icon) icon.classList.remove('hidden');
        if (spinner) spinner.classList.add('hidden');
        btn.disabled = false;
    }

    // Shows clear, actionable guidance when geolocation is denied by the browser
    function showLocationDeniedGuidance() {
        const { btn } = getElements();
        if (!btn) return;

        applyDeniedState('Location blocked by browser');

        // Find or create the guidance message next to the button container
        const container = btn.parentElement;
        if (!container) return;

        let guidance = container.querySelector('#near-me-denied-guidance');
        if (!guidance) {
            guidance = document.createElement('div');
            guidance.id = 'near-me-denied-guidance';
            guidance.className = 'text-[10px] text-text-sub mt-1 w-full max-w-[220px] leading-tight';
            container.appendChild(guidance);
        }

        guidance.innerHTML = 
            'To see listings near you, allow location for agbalumo.com in your browser settings ' +
            '(click the lock icon in the address bar), then tap <strong>NEAR ME</strong> again.';
        guidance.style.display = 'block';
    }

    function hideLocationDeniedGuidance() {
        const container = document.getElementById('near-me-btn')?.parentElement;
        if (!container) return;

        const guidance = container.querySelector('#near-me-denied-guidance');
        if (guidance) guidance.style.display = 'none';
    }

    // Core geolocation request helper
    function requestGeolocation(onSuccess, onError) {
        if (!("geolocation" in navigator)) {
            if (onError) onError(new Error('Geolocation not supported'));
            return;
        }

        navigator.geolocation.getCurrentPosition(
            (position) => {
                const lat = position.coords.latitude;
                const lng = position.coords.longitude;

                sessionStorage.setItem('agbalumo_lat', lat);
                sessionStorage.setItem('agbalumo_lng', lng);

                if (onSuccess) onSuccess(lat, lng);
            },
            (error) => {
                console.warn('Geolocation request failed:', error);
                if (onError) onError(error);
            },
            { enableHighAccuracy: false, timeout: 10000, maximumAge: 60000 }
        );
    }

    function initNearMe() {
        const { btn } = getElements();
        if (!btn) return;

        // Seed from URL if present (e.g. server-side or shared link)
        const urlParams = new URLSearchParams(window.location.search);
        const urlLat = urlParams.get('lat');
        const urlLng = urlParams.get('lng');
        if (urlLat && urlLng) {
            sessionStorage.setItem('agbalumo_lat', urlLat);
            sessionStorage.setItem('agbalumo_lng', urlLng);
        }

        // On initial load: immediately attempt to get location for Ada's "find food near her fast" intent.
        // This triggers the browser's native permission prompt on first visit.
        const existingLat = sessionStorage.getItem('agbalumo_lat');
        const existingLng = sessionStorage.getItem('agbalumo_lng');

        if (!existingLat || !existingLng) {
            // Option C: Proactive check using Permissions API
            if (navigator.permissions && navigator.permissions.query) {
                navigator.permissions.query({ name: 'geolocation' }).then(result => {
                    if (result.state === 'denied') {
                        showLocationDeniedGuidance();
                        return;
                    }
                    // Not denied yet — attempt the request (may trigger native prompt)
                    attemptInitialGeolocation();
                }).catch(() => {
                    // Permissions API unavailable — fall back to direct attempt
                    attemptInitialGeolocation();
                });
            } else {
                attemptInitialGeolocation();
            }
        } else {
            applyActiveState();
        }

        function attemptInitialGeolocation() {
            requestGeolocation(
                (lat, lng) => {
                    applyActiveState();
                    hideLocationDeniedGuidance();
                    if (window.htmx) {
                        const values = {
                            lat: lat,
                            lng: lng,
                            radius: window.filterState?.radius || '10'
                        };
                        htmx.ajax('GET', '/listings/fragment', {
                            values: values,
                            target: '#listings-container',
                            swap: 'innerHTML'
                        });
                    }
                },
                (error) => {
                    // Silent on initial load — user can still click NEAR ME later.
                    applyDefaultState();
                }
            );
        }

        btn.addEventListener('click', function() {
            const { btn: clickBtn } = getElements();
            if (!clickBtn) return;

            const lat = sessionStorage.getItem('agbalumo_lat');
            const lng = sessionStorage.getItem('agbalumo_lng');

            // If we already have location → toggle off (current good behavior)
            if (lat && lng) {
                sessionStorage.removeItem('agbalumo_lat');
                sessionStorage.removeItem('agbalumo_lng');
                applyDefaultState();
                hideLocationDeniedGuidance();

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

            // No location yet: attempt geolocation again on explicit NEAR ME click.
            // This is the user's explicit request to try getting location.
            // If previously denied, the browser will return error (we cannot force re-prompt).
            clickBtn.disabled = true;

            requestGeolocation(
                (lat, lng) => {
                    applyActiveState();
                    hideLocationDeniedGuidance();

                    if (window.htmx) {
                        const values = {
                            lat: lat,
                            lng: lng,
                            radius: window.filterState?.radius || '10'
                        };
                        htmx.ajax('GET', '/listings/fragment', {
                            values: values,
                            target: '#listings-container',
                            swap: 'innerHTML'
                        });
                    }
                },
                (error) => {
                    const isDenied = error.code === 1 || (error instanceof GeolocationPositionError && error.code === error.PERMISSION_DENIED);

                    if (isDenied) {
                        showLocationDeniedGuidance();
                    } else {
                        applyDeniedState('Unable to get location');
                    }
                    clickBtn.disabled = false;
                }
            );
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
