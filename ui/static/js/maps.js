// Google Maps lazy loading and utilities

let googleMapsLoaded = false;

function loadGoogleMapsApi(apiKey) {
    if (googleMapsLoaded || !apiKey) return;

    const script = document.createElement('script');
    script.src = `https://maps.googleapis.com/maps/api/js?key=${apiKey}&libraries=places&callback=initGoogleMaps`;
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
    googleMapsLoaded = true;
}

// Google Maps Init (Global scope needed for callback)
window.initGoogleMaps = function () {
    const inputs = document.querySelectorAll('[name="address"][data-google-maps-key]');
    if (inputs.length === 0) return;

    if (typeof google === 'undefined' || !google.maps || !google.maps.places) return;

    inputs.forEach(input => {
        // Prevent double-initialization
        if (input.dataset.autocompleteInitialized) return;

        const options = {
            fields: ["address_components", "geometry", "formatted_address"],
            types: ["address"],
        };
        const autocomplete = new google.maps.places.Autocomplete(input, options);
        autocomplete.addListener("place_changed", () => {
            const place = autocomplete.getPlace();
            if (!place || !place.address_components) return;

            let city = "";
            for (const component of place.address_components) {
                const types = component.types;
                if (types.includes("locality")) {
                    city = component.long_name;
                    break;
                } else if (types.includes("sublocality_level_1") || types.includes("sublocality")) {
                    city = component.long_name;
                } else if (!city && (types.includes("postal_town") || types.includes("administrative_area_level_2") || types.includes("neighborhood"))) {
                    city = component.long_name;
                }
            }

            // Find companion city input based on ID/Structure
            const id = input.id;
            const cityInputId = id.replace('address-input', 'city-input');
            const cityInput = document.getElementById(cityInputId);

            if (cityInput) {
                cityInput.value = city || "Unknown";
            }

            if (place.formatted_address) {
                input.value = place.formatted_address;
            }
        });
        input.dataset.autocompleteInitialized = "true";
    });
}

// Load Google Maps when create or edit listing modal opens
function setupGoogleMapsLazyLoad() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-id="create-listing-modal"], [hx-get*="/edit"]');
        if (btn) {
            // Short delay to ensure modal is in DOM or being loaded
            setTimeout(() => {
                const input = document.querySelector('[name="address"][data-google-maps-key]');
                if (input) {
                    const apiKey = input.dataset.googleMapsKey;
                    if (apiKey) {
                        loadGoogleMapsApi(apiKey);
                    }
                }
            }, 500);
        }
    });

    // Also look for edit modals already in DOM (HTMX swapped)
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        const elt = evt.detail.elt;
        const input = elt.querySelector ? elt.querySelector('[name="address"][data-google-maps-key]') : null;
        if (input) {
            const apiKey = input.dataset.googleMapsKey;
            if (apiKey) {
                loadGoogleMapsApi(apiKey);
            }
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    setupGoogleMapsLazyLoad();
});
