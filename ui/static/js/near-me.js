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

        // Persist active state on load
        const lat = sessionStorage.getItem('agbalumo_lat');
        const lng = sessionStorage.getItem('agbalumo_lng');
        if (lat && lng) {
            applyActiveState();
        }

        btn.addEventListener('click', function() {
            const { btn: clickBtn, icon, spinner, text } = getElements();
            if (!clickBtn) return;

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

                        // Visual feedback transition
                        const container = document.getElementById('listings-container');
                        if (container) {
                            container.style.opacity = '0.3';
                        }

                        if (window.htmx) {
                            const urlParams = new URLSearchParams(window.location.search);
                            const values = {
                                lat: lat,
                                lng: lng,
                                radius: 10
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
                        let errMsg = 'Error: Unavailable';
                        if (error.code === error.PERMISSION_DENIED) {
                            errMsg = 'Error: Denied';
                        }
                        if (text) text.textContent = errMsg;
                        if (spinner) spinner.classList.add('hidden');

                        setTimeout(function() {
                            applyDefaultState();
                        }, 2000);
                    },
                    { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
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
                        evt.detail.parameters['radius'] = '10';
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
