// agbalumo Main Application Logic

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
        if (event.target.matches('#create-listing-modal select[name="type"], #create-listing-modal input[name="type"]')) {
            toggleListingFields(event.target);
        }
    });

    // If the modal is already present or inserted, we might need to init.
    // MutationObserver could be used, or just running check on interactions.
    // For now, let's also look for the element on load.
    const typeSelect = document.querySelector('#create-listing-modal select[name="type"], #create-listing-modal input[name="type"]');
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

        // Auto-open and init edit modals inserted via HTMX
        function initEditDialog(d) {
            if (d.dataset.autoOpen === 'true' && !d.open) {
                d.showModal();
                const listingId = d.dataset.listingId;
                if (listingId) {
                    if (window.initEditTypeToggle) window.initEditTypeToggle(listingId);
                    if (window.initEditMaps) window.initEditMaps(listingId);
                    if (window.initEditImagePreview) window.initEditImagePreview(listingId);
                }
            }
        }

        // Check if the swapped element is the dialog itself
        if (elt.tagName === 'DIALOG' && elt.id.startsWith('edit-listing-modal-')) {
            initEditDialog(elt);
        } else {
            // Or if the dialog is contained within the swapped element
            const dialogs = elt.querySelectorAll ? elt.querySelectorAll('dialog[id^="edit-listing-modal-"]') : [];
            dialogs.forEach(d => initEditDialog(d));
        }

        // Init create image preview when create modal is swapped in
        if (elt.id === 'create-listing-modal' || (elt.querySelectorAll && elt.querySelector('#create-listing-modal'))) {
            if (window.initCreateImagePreview) window.initCreateImagePreview();
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

    // Initialize custom dropdowns on HTMX swap
    document.body.addEventListener('htmx:afterSwap', (evt) => {
        if (typeof initCustomDropdownsActiveState === 'function') {
            initCustomDropdownsActiveState(evt.detail.elt);
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

    document.body.addEventListener('htmx:afterSwap', (evt) => {
        evt.detail.elt.querySelectorAll('.dynamic-bg').forEach(el => {
            observer.observe(el);
        });
    });
}

// 9. Featured Carousel
function setupFeaturedCarousel() {
    const carousel = document.getElementById('featured-carousel');
    if (!carousel) return;

    const slides = carousel.querySelectorAll('.carousel-slide');
    const dots = carousel.querySelectorAll('.carousel-dot');
    const prevBtn = carousel.querySelector('.carousel-prev');
    const nextBtn = carousel.querySelector('.carousel-next');

    if (slides.length <= 1) return;

    let currentIndex = 0;
    let interval = null;
    let isPaused = false;

    function goToSlide(index) {
        slides.forEach((slide, i) => {
            slide.classList.toggle('opacity-100', i === index);
            slide.classList.toggle('z-10', i === index);
            slide.classList.toggle('opacity-0', i !== index);
            slide.classList.toggle('z-0', i !== index);
        });

        dots.forEach((dot, i) => {
            dot.classList.toggle('bg-white', i === index);
            dot.classList.toggle('w-6', i === index);
            dot.classList.toggle('md:w-8', i === index);
            dot.classList.toggle('bg-white/40', i !== index);
            dot.classList.toggle('w-2', i !== index);
            dot.classList.toggle('md:w-3', i !== index);
        });

        currentIndex = index;
    }

    function nextSlide() {
        goToSlide((currentIndex + 1) % slides.length);
    }

    function prevSlide() {
        goToSlide((currentIndex - 1 + slides.length) % slides.length);
    }

    function startAutoplay() {
        if (interval) clearInterval(interval);
        const delay = 5000 + Math.random() * 3000;
        interval = setInterval(() => {
            if (!isPaused) nextSlide();
        }, delay);
    }

    function pauseAutoplay() {
        isPaused = true;
    }

    function resumeAutoplay() {
        isPaused = false;
    }

    if (prevBtn) {
        prevBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            prevSlide();
            startAutoplay();
        });
    }

    if (nextBtn) {
        nextBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            nextSlide();
            startAutoplay();
        });
    }

    dots.forEach((dot, i) => {
        dot.addEventListener('click', (e) => {
            e.stopPropagation();
            goToSlide(i);
            startAutoplay();
        });
    });

    carousel.addEventListener('mouseenter', pauseAutoplay);
    carousel.addEventListener('mouseleave', resumeAutoplay);
    carousel.addEventListener('focusin', pauseAutoplay);
    carousel.addEventListener('focusout', resumeAutoplay);

    if (carousel.dataset.autoplay === 'true') {
        startAutoplay();
    }
}

// 10. Create Image Preview Init (replaces inline script)
function setupCreateImagePreviewInit() {
    // Init on page load if create modal already exists in DOM
    if (document.getElementById('create-listing-modal') && window.initCreateImagePreview) {
        window.initCreateImagePreview();
    }

    // Also init when create modal is opened via button click
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-id="create-listing-modal"]');
        if (btn) {
            setTimeout(() => {
                if (window.initCreateImagePreview) window.initCreateImagePreview();
            }, 50);
        }
    });
}

// 11. Filter Toggle Logic
function setupFilterToggle() {
    const activeClasses = ['bg-earth-ochre', 'text-earth-dark'];
    const inactiveClasses = ['bg-white/10', 'text-white'];

    const updateButtonStates = (isFilterOpen) => {
        const filtersBtns = document.querySelectorAll('.filters-btn');
        const searchBtns = document.querySelectorAll('.search-btn');

        if (isFilterOpen) {
            filtersBtns.forEach(btn => {
                btn.classList.add(...activeClasses);
                btn.classList.remove(...inactiveClasses);
            });
            searchBtns.forEach(btn => {
                btn.classList.remove(...activeClasses);
                btn.classList.add(...inactiveClasses);
            });
        } else {
            filtersBtns.forEach(btn => {
                btn.classList.remove(...activeClasses);
                btn.classList.add(...inactiveClasses);
            });
            searchBtns.forEach(btn => {
                btn.classList.add(...activeClasses);
                btn.classList.remove(...inactiveClasses);
            });
        }
    };

    document.addEventListener('click', (e) => {
        const filtersBtn = e.target.closest('[data-action="toggle-filters"]');
        if (filtersBtn) {
            const panel = document.getElementById('filter-dropdown-panel');
            if (panel) {
                const isWillBeOpen = panel.classList.contains('hidden');
                panel.classList.toggle('hidden');
                updateButtonStates(isWillBeOpen);

                // Mobile Bottom Sheet & Scroll Lock
                if (window.innerWidth < 768) {
                    const bottomNav = document.getElementById('mobile-bottom-nav');
                    let overlay = document.getElementById('mobile-filter-overlay');
                    if (!overlay) {
                        overlay = document.createElement('div');
                        overlay.id = 'mobile-filter-overlay';
                        // z-[105] is between header (110) and body (0)
                        overlay.className = 'fixed inset-0 bg-earth-dark/70 z-[105] hidden transition-opacity duration-300 backdrop-blur-[2px]';
                        document.body.appendChild(overlay);
                        overlay.onclick = () => {
                            if (panel) panel.classList.add('hidden');
                            overlay.classList.add('hidden');
                            if (bottomNav) bottomNav.classList.remove('nav-hidden-js');
                            document.body.style.overflow = '';
                            updateButtonStates(false);
                        };
                    }

                    if (isWillBeOpen) {
                        overlay.classList.remove('hidden');
                        if (bottomNav) bottomNav.classList.add('nav-hidden-js');
                        document.body.style.overflow = 'hidden';
                    } else {
                        overlay.classList.add('hidden');
                        if (bottomNav) bottomNav.classList.remove('nav-hidden-js');
                        document.body.style.overflow = '';
                    }
                }
            }
            return;
        }

        const htmxTriggerBtn = e.target.closest('[data-filter-type]');
        const searchBtn = e.target.closest('.search-btn');
        if (searchBtn || htmxTriggerBtn) {
            const panel = document.getElementById('filter-dropdown-panel');
            const overlay = document.getElementById('mobile-filter-overlay');
            const bottomNav = document.getElementById('mobile-bottom-nav');
            if (panel && !panel.classList.contains('hidden')) {
                panel.classList.add('hidden');
                if (overlay) overlay.classList.add('hidden');
                if (bottomNav) bottomNav.classList.remove('nav-hidden-js');
                if (window.innerWidth < 768) {
                    document.body.style.overflow = '';
                }
            }
            if (searchBtn) {
                updateButtonStates(false);
                return;
            }
        }

        // Filter chip click handler
        const chip = e.target.closest('[data-filter-type]');
        if (chip) {
            const filterType = chip.dataset.filterType;
            const filterValue = chip.dataset.filterValue;
            const container = chip.parentElement;

            // Update active chip styles within the same group
            if (container) {
                container.querySelectorAll('button').forEach(b => {
                    b.classList.remove('bg-earth-ochre', 'text-white');
                    b.classList.add('bg-earth-dark/5', 'text-earth-dark');
                });
                chip.classList.add('bg-earth-ochre', 'text-white');
                chip.classList.remove('bg-earth-dark/5', 'text-earth-dark');
            }

            // Build HTMX request URL
            let url = '/listings/fragment';
            if (filterType === 'category' && filterValue) {
                url += '?type=' + encodeURIComponent(filterValue);
            } else if (filterType === 'location' && filterValue) {
                url += '?q=' + encodeURIComponent(filterValue);
            }
            htmx.ajax('GET', url, '#listings-container');
        }
    });
}

// 12. Custom Dropdown Logic
function initCustomDropdownsActiveState(root = document) {
    // Initialize active states based on hidden inputs
    root.querySelectorAll('.custom-dropdown').forEach(container => {
        const input = container.querySelector('input[type="hidden"]');
        if (input && input.value) {
            const activeBtn = container.querySelector(`[data-dropdown-value="${input.value}"]`);
            if (activeBtn) {
                container.querySelectorAll('[data-dropdown-value]').forEach(btn => {
                    btn.classList.remove('bg-earth-ochre', 'text-white');
                    btn.classList.add('bg-earth-dark/5', 'text-earth-dark');
                });
                activeBtn.classList.remove('bg-earth-dark/5', 'text-earth-dark');
                activeBtn.classList.add('bg-earth-ochre', 'text-white');
            }
        }
    });
}

function setupCustomDropdowns() {
    initCustomDropdownsActiveState();
    
    

    document.addEventListener('click', (e) => {
        // Handle clicking a dropdown toggle button
        const toggleBtn = e.target.closest('[data-dropdown-toggle]');
        if (toggleBtn) {
            const container = toggleBtn.closest('.custom-dropdown');
            if (!container) return;
            
            const menu = container.querySelector('.dropdown-menu');
            if (menu) {
                // Close other open dropdowns first
                document.querySelectorAll('.custom-dropdown .dropdown-menu').forEach(m => {
                    if (m !== menu) m.classList.add('hidden');
                });
                
                menu.classList.toggle('hidden');
                
                // Adjust z-index of all containers to ensure the active one sits on top
                document.querySelectorAll('.custom-dropdown').forEach(dropdown => {
                    dropdown.style.zIndex = dropdown === container ? '50' : '10';
                });
            }
            return;
        }

        // Handle clicking an option inside the dropdown menu
        const optionBtn = e.target.closest('[data-dropdown-value]');
        if (optionBtn) {
            const container = optionBtn.closest('.custom-dropdown');
            if (!container) return;

            const value = optionBtn.dataset.dropdownValue;
            let label = optionBtn.textContent.trim();
            
            // Remove emoji flags from label for display if we want cleaner text,
            // but for now we'll just use the textContent directly.
            
            const input = container.querySelector('input[type="hidden"]');
            const display = container.querySelector('.dropdown-display');
            const menu = container.querySelector('.dropdown-menu');

            if (input && input.value !== value) {
                input.value = value;
                // Dispatch change event so other scripts (like toggleListingFields) catch it
                input.dispatchEvent(new Event('change', { bubbles: true }));
            }

            if (display) {
                display.textContent = label;
            }

            // Update active state visuals
            container.querySelectorAll('[data-dropdown-value]').forEach(btn => {
                btn.classList.remove('bg-earth-ochre', 'text-white');
                btn.classList.add('bg-earth-dark/5', 'text-earth-dark');
            });
            optionBtn.classList.remove('bg-earth-dark/5', 'text-earth-dark');
            optionBtn.classList.add('bg-earth-ochre', 'text-white');

            if (menu) {
                menu.classList.add('hidden');
                
                // Reset z-index
                container.style.zIndex = '10';
            }
            return;
        }

        // Click outside closes all dropdowns
        if (!e.target.closest('.custom-dropdown')) {
            document.querySelectorAll('.custom-dropdown .dropdown-menu').forEach(m => {
                m.classList.add('hidden');
                const container = m.closest('.custom-dropdown');
                if (container) container.style.zIndex = '10';
            });
        }
    });
}

const originalInit = initApp;
initApp = function () {
    originalInit();
    setupAuthActions();
    setupModalActions();
    setupHtmxIntegration();
    setupCsrf();
    setupDynamicBackgrounds();
    setupGoogleMapsLazyLoad();
    setupFeaturedCarousel();
    setupCreateImagePreviewInit();
    setupFilterToggle();
    setupCustomDropdowns();
};
