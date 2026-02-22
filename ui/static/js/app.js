// Agbalumo Main Application Logic

document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

function initApp() {
    setupModalClosing();
    setupFilterButtons();
    setupMobileBottomNav();
    // Re-initialize logic when HTMX swaps content if necessary
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        // specific re-init if needed
    });

    // Listing Modal Logic (Event Delegation because modal might not exist yet)
    setupListingModalDelegation();
}

// 0. Mobile Bottom Nav Scroll-Aware Show/Hide
function setupMobileBottomNav() {
    const nav = document.getElementById('mobile-bottom-nav');
    if (!nav) return;

    // Use the main scrollable element, not window
    const scrollContainer = document.querySelector('main');
    if (!scrollContainer) return;

    let lastScrollY = 0;
    let scrollTimeout;

    scrollContainer.addEventListener('scroll', () => {
        const currentScrollY = scrollContainer.scrollTop;
        const isScrollingDown = currentScrollY > lastScrollY && currentScrollY > 60;

        if (isScrollingDown) {
            nav.classList.add('nav-hidden');
        } else {
            nav.classList.remove('nav-hidden');
        }

        lastScrollY = currentScrollY;

        // Always show after user stops scrolling
        clearTimeout(scrollTimeout);
        scrollTimeout = setTimeout(() => {
            nav.classList.remove('nav-hidden');
        }, 1500);
    }, { passive: true });
}

// 1. Modal Closing Logic
function setupModalClosing() {
    document.addEventListener('click', (event) => {
        if (event.target.tagName === 'DIALOG') {
            const dialog = event.target;
            const rect = dialog.getBoundingClientRect();
            const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height &&
                rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                dialog.close();
                // Remove if it's a dynamic modal
                if (dialog.id.startsWith('detail-modal-') || dialog.id.startsWith('edit-modal-')) {
                    dialog.remove();
                }
            }
        }
    });
}

// 2. Filter Buttons Logic
function setupFilterButtons() {
    document.addEventListener('click', (event) => {
        const btn = event.target.closest('#filter-container button');
        if (btn) {
            const container = document.getElementById('filter-container');
            const buttons = container.querySelectorAll('button');
            const activeState = "chip-fruit flex h-11 shrink-0 items-center justify-center gap-x-2 rounded-full bg-stone-900 text-white dark:bg-white dark:text-stone-900 px-4 shadow-sm border border-transparent transition-transform active:scale-95 text-xs font-bold uppercase";
            const inactiveState = "chip-fruit flex h-11 shrink-0 items-center justify-center gap-x-2 rounded-full bg-white dark:bg-surface-dark border border-stone-200 dark:border-stone-700 px-4 transition-transform active:scale-95 hover:bg-stone-50 text-text-main dark:text-stone-200 text-xs font-semibold uppercase";

            buttons.forEach(b => {
                b.className = inactiveState;
                const icon = b.querySelector('.material-symbols-outlined');
                if (icon) {
                    if (!icon.classList.contains('text-text-sub')) {
                        icon.classList.add('text-text-sub');
                    }
                }
            });

            btn.className = activeState;
            const icon = btn.querySelector('.material-symbols-outlined');
            if (icon) {
                icon.classList.remove('text-text-sub');
            }
        }
    }, true);
}

// 3. Create Listing Modal Logic
// We attach this to document because the modal is inserted dynamically or exists statically but hidden.
// Since it's a dialog, we can check mostly when it opens or just delegate change events.
function setupListingModalDelegation() {
    document.addEventListener('change', (event) => {
        if (event.target.matches('#create-listing-modal select[name="type"]')) {
            toggleListingFields(event.target);
        }
    });

    // If the modal is already present or inserted, we might need to init.
    // MutationObserver could be used, or just running check on interactions.
    // For now, let's also look for the element on load.
    const typeSelect = document.querySelector('#create-listing-modal select[name="type"]');
    if (typeSelect) {
        toggleListingFields(typeSelect);
    }
}

function toggleListingFields(typeSelect) {
    const modal = typeSelect.closest('dialog');
    if (!modal) return;

    const eventSection = modal.querySelector('#event-dates-section');
    const jobSection = modal.querySelector('#job-fields-section');
    const imageSection = modal.querySelector('#image-upload-section');
    const hoursSection = modal.querySelector('#hours-section');
    const locationLabel = modal.querySelector('#location-label');
    const descriptionLabel = modal.querySelector('#description-label');
    const addressInput = modal.querySelector('#create-address-input');

    const val = typeSelect.value;

    // Event Logic
    if (eventSection) {
        if (val === 'Event') {
            eventSection.classList.remove('hidden');
            eventSection.querySelectorAll('input').forEach(i => i.required = true);
        } else {
            eventSection.classList.add('hidden');
            eventSection.querySelectorAll('input').forEach(i => {
                i.required = false;
                i.value = '';
            });
        }
    }

    // Job Logic
    if (jobSection) {
        if (val === 'Job') {
            jobSection.classList.remove('hidden');
            jobSection.querySelectorAll('input, textarea').forEach(i => {
                if (i.name === 'company' || i.name === 'skills' || i.name === 'job_start_date' || i.name === 'pay_range' || i.name === 'job_apply_url') {
                    i.required = true;
                }
            });

            // Hide Image Upload for Jobs
            if (imageSection) imageSection.classList.add('hidden');
        } else {
            jobSection.classList.add('hidden');
            jobSection.querySelectorAll('input, textarea').forEach(i => {
                i.required = false;
                i.value = '';
            });

            // Show Image Upload for others
            if (imageSection) imageSection.classList.remove('hidden');
        }
    }

    // Address & Label Logic
    if (addressInput) {
        if (val === 'Job') {
            addressInput.required = true;
            addressInput.placeholder = "City, Country or Address";
            if (locationLabel) locationLabel.textContent = "Location";
            if (descriptionLabel) descriptionLabel.textContent = "Job Description";
        } else {
            if (descriptionLabel) descriptionLabel.textContent = "Description";
            // Address is optional for Service, Request, Product, and Event
            if (val === 'Service' || val === 'Request' || val === 'Product' || val === 'Event') {
                addressInput.required = false;
                addressInput.placeholder = "Address (Optional)";
                if (locationLabel) locationLabel.textContent = "Address (Optional)";
            } else {
                // Business, Food
                addressInput.required = true;
                addressInput.placeholder = "Start typing address...";
                if (locationLabel) locationLabel.textContent = "Address (Validated)";
            }
        }
    }

    // Hours of Operation Logic
    if (hoursSection) {
        const hoursInput = hoursSection.querySelector('input');
        // Allowed: Business, Service, Food
        // Explicitly check for allowed types
        if (val === 'Business' || val === 'Service' || val === 'Food') {
            hoursSection.classList.remove('hidden');
            hoursSection.style.display = ''; // Clear any inline styles
            if (hoursInput) hoursInput.disabled = false;
        } else {
            // Product, Event, Job, Request
            hoursSection.classList.add('hidden');
            hoursSection.style.display = ''; // Clear any inline styles
            if (hoursInput) {
                hoursInput.value = '';
                hoursInput.disabled = true;
            }
        }
    }
}

// Google Maps lazy loading
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
    const input = document.getElementById('create-address-input');
    if (!input) return;

    if (!google || !google.maps || !google.maps.places) return;

    const options = {
        fields: ["address_components", "geometry", "formatted_address"],
        types: ["address"],
    };
    const autocomplete = new google.maps.places.Autocomplete(input, options);
    autocomplete.addListener("place_changed", () => {
        const place = autocomplete.getPlace();
        let city = "";
        // Find city/locality
        if (place.address_components) {
            for (const component of place.address_components) {
                const types = component.types;
                if (types.includes("locality")) {
                    city = component.long_name;
                    break;
                }
                if (types.includes("postal_town")) {
                    city = component.long_name;
                }
                if (!city && types.includes("administrative_area_level_2")) {
                    city = component.long_name;
                }
            }
        }

        const cityInput = document.getElementById('create-city-input');
        if (cityInput) {
            if (city) {
                cityInput.value = city;
            } else {
                // Fallback
                cityInput.value = "Unknown";
            }
        }

        if (place.formatted_address) {
            input.value = place.formatted_address;
        }
    });
}

// Load Google Maps when create-listing modal opens
function setupGoogleMapsLazyLoad() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-id="create-listing-modal"]');
        if (btn) {
            const input = document.getElementById('create-address-input');
            if (input) {
                const apiKey = input.dataset.googleMapsKey;
                if (apiKey) {
                    loadGoogleMapsApi(apiKey);
                }
            }
        }
    });
}
// 4. Auth Action Logic (Replaces inline onclicks)
function setupAuthActions() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-auth-action="modal"]');
        if (btn) {
            const isAuthenticated = btn.dataset.userAuthenticated === 'true';
            if (isAuthenticated) {
                const modalId = btn.dataset.modalId;
                const modal = document.getElementById(modalId);
                if (modal) modal.showModal();
            } else {
                window.location.href = '/auth/google/login';
            }
        }
    });
}

// 5. Generic Modal Actions (Close, Submit)
function setupModalActions() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-action="close"]');
        if (btn) {
            // Find closest dialog or use ID
            const modal = btn.closest('dialog') || document.getElementById(btn.dataset.modalId);
            if (modal) {
                modal.close();
                // Reset form fields in static create dialogs so next open is blank
                const form = modal.querySelector('form');
                if (form) form.reset();
                // Remove dynamic modals (detail/edit) from DOM
                if (modal.id.startsWith('detail-modal-') || modal.id.startsWith('edit-modal-')) {
                    modal.remove();
                }
            }
        }
    });

    document.addEventListener('submit', (e) => {
        const form = e.target;
        if (form.hasAttribute('data-modal-action-submit-close')) {
            const modal = form.closest('dialog');
            if (modal) {
                // Wait for HTMX to finish, then reset + close
                modal.addEventListener('htmx:afterRequest', () => {
                    modal.close();
                    form.reset();
                }, { once: true });
            }
        }
    });
}

// 6. HTMX Integration for Modals
function setupHtmxIntegration() {
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        // htmx:afterSwap provides the swapped element in evt.detail.elt
        const elt = evt.detail.elt;
        if (!elt) return;

        // Check if the swapped element is the dialog itself
        if (elt.tagName === 'DIALOG' && elt.id.startsWith('edit-listing-modal-')) {
            if (!elt.open) elt.showModal();
        } else {
            // Or if the dialog is contained within the swapped element
            const dialogs = elt.querySelectorAll ? elt.querySelectorAll('dialog[id^="edit-listing-modal-"]') : [];
            dialogs.forEach(d => {
                if (!d.open) d.showModal();
            });
        }
    });

    // Auto-close edit modal after successful save
    document.body.addEventListener('htmx:afterRequest', (evt) => {
        const form = evt.detail.elt;
        if (!form || !form.id || !form.id.startsWith('edit-form-')) return;
        if (!evt.detail.successful) return;

        const dialog = form.closest('dialog');
        if (dialog) {
            dialog.close();
            dialog.remove();
        }
    });
}

// 7. CSRF Token Injection for HTMX
function setupCsrf() {
    document.body.addEventListener('htmx:configRequest', (evt) => {
        const token = document.querySelector('meta[name="csrf-token"]').getAttribute('content');
        evt.detail.headers['X-CSRF-Token'] = token;
    });
}

// 8. Dynamic Background Images (CSP Safe)
function setupDynamicBackgrounds() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const el = entry.target;
                const url = el.dataset.bgImage;
                if (url) {
                    el.style.backgroundImage = `url('${url}')`;
                    el.style.backgroundSize = 'cover';
                    el.style.backgroundPosition = 'center';
                }
                observer.unobserve(el);
            }
        });
    });

    document.querySelectorAll('.dynamic-bg').forEach(el => {
        observer.observe(el);
    });

    // Also handle HTMX swaps
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        evt.detail.elt.querySelectorAll('.dynamic-bg').forEach(el => {
            observer.observe(el);
        });
    });
}

// Add to init
const originalInit = initApp;
initApp = function () {
    originalInit();
    setupAuthActions();
    setupModalActions();
    setupHtmxIntegration();
    setupCsrf();
    setupDynamicBackgrounds();
    setupGoogleMapsLazyLoad();
};
