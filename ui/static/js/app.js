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
    setupDynamicBackgrounds();
    setupFeaturedCarousel();
    setupFilterToggle();
    setupCustomDropdowns();
};
