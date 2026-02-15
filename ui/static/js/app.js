// Agbalumo Main Application Logic

document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

function initApp() {
    setupModalClosing();
    setupFilterButtons();
    // Re-initialize logic when HTMX swaps content if necessary
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        // specific re-init if needed
    });

    // Listing Modal Logic (Event Delegation because modal might not exist yet)
    setupListingModalDelegation();
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
            const activeState = "flex h-8 shrink-0 items-center justify-center gap-x-2 rounded-full bg-stone-900 text-white dark:bg-white dark:text-stone-900 px-4 shadow-sm border border-transparent transition-transform active:scale-95 text-xs font-bold uppercase";
            const inactiveState = "flex h-8 shrink-0 items-center justify-center gap-x-2 rounded-full bg-white dark:bg-surface-dark border border-stone-200 dark:border-stone-700 px-4 transition-transform active:scale-95 hover:bg-stone-50 text-text-main dark:text-stone-200 text-xs font-semibold uppercase";

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
            if (val === 'Service' || val === 'Request') {
                addressInput.required = false;
                addressInput.placeholder = "Address (Optional)";
                if (locationLabel) locationLabel.textContent = "Address (Optional)";
            } else {
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
        if (val === 'Business' || val === 'Service' || val === 'Food') {
            hoursSection.classList.remove('hidden');
            if (hoursInput) hoursInput.disabled = false;
        } else {
            hoursSection.classList.add('hidden');
            if (hoursInput) {
                hoursInput.value = '';
                hoursInput.disabled = true;
            }
        }
    }
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
                modal.close();
            }
        }
    });
}

// 6. HTMX Integration for Modals
function setupHtmxIntegration() {
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        // If content is a dialog (or contains one), check if it needs opening
        // Specifically for edit modals which we stripped the <script> from.
        // We can identify them by ID prefix 'edit-listing-modal-'
        const target = evt.target;
        // HTMX swap target might be the element itself or parent.
        // If we swapped outerHTML of a dialog, target is the new dialog? 
        // No, target is the element designated by hx-target.
        // But let's check the added nodes or look for the dialog in DOM.

        // A simpler approach: if the swapped content is a dialog with known prefix, show it.
        // evt.detail.elt is the swapped element.
        const elt = evt.detail.elt;
        if (elt && elt.tagName === 'DIALOG' && elt.id.startsWith('edit-listing-modal-')) {
            elt.showModal();
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
};
