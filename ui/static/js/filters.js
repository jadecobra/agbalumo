/**
 * Filters Component Logic
 * Robust discovery filtering for the agbalumo platform.
 * CSP-compliant: Standard DOM listeners, avoid inline Alpine.
 */

// Global state for HTMX hx-vals
if (!window.filterState) {
    const urlParams = new URLSearchParams(window.location.search);
    window.filterState = {
        type: urlParams.get('type') || 'Food',
        city: urlParams.get('city') || '',
        radius: urlParams.get('radius') || '25'
    };
}

// Global function to toggle the filter panel and handle all accessibility properties
function toggleFilterPanel(show) {
    const panel = document.getElementById('filter-dropdown-panel');
    if (!panel) return;

    if (show === undefined) {
        show = panel.classList.contains('hidden');
    }
    
    const triggers = document.querySelectorAll('[data-testid^="ag-home-filters-toggle"], [data-testid="ag-home-filters-close"]');
    
    if (show) {
        panel.classList.remove('hidden');
        
        // Update aria states on trigger elements
        triggers.forEach(trigger => trigger.setAttribute('aria-expanded', 'true'));

        // Viewport-aware inert focus trap on mobile
        if (window.innerWidth < 768) {
            const mainContent = document.getElementById('main-content');
            if (mainContent) mainContent.setAttribute('inert', '');
            const mobileNav = document.getElementById('mobile-bottom-nav');
            if (mobileNav) mobileNav.setAttribute('inert', '');
        }

        // Focus search input on open
        const searchInput = document.getElementById('search-nav') || document.getElementById('search');
        if (searchInput) setTimeout(() => searchInput.focus(), 100);
    } else {
        panel.classList.add('hidden');

        // Update aria states on trigger elements
        triggers.forEach(trigger => trigger.setAttribute('aria-expanded', 'false'));

        // Clean up inert attributes on close
        const mainContent = document.getElementById('main-content');
        if (mainContent) mainContent.removeAttribute('inert');
        const mobileNav = document.getElementById('mobile-bottom-nav');
        if (mobileNav) mobileNav.removeAttribute('inert');
    }
}

function setupFilterToggle() {
    if (window._filterToggleInitialized) return;
    window._filterToggleInitialized = true;

    document.addEventListener('click', (e) => {
        const toggle = e.target.closest('[data-testid^="ag-home-filters-toggle"]');
        if (toggle) {
            e.stopPropagation();
            toggleFilterPanel();
            return;
        }

        const closeBtn = e.target.closest('[data-testid="ag-home-filters-close"]');
        if (closeBtn) {
            e.stopPropagation();
            toggleFilterPanel(false);
            return;
        }

        const panel = document.getElementById('filter-dropdown-panel');
        if (panel && !panel.contains(e.target) && !panel.classList.contains('hidden')) {
            toggleFilterPanel(false);
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') toggleFilterPanel(false);
    });
}

function setupFilterButtons() {
    if (window._filterButtonsInitialized) return;
    window._filterButtonsInitialized = true;

    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-category-name]');
        if (!btn) return;

        const category = btn.getAttribute('data-category-name') || '';
        
        // Update global state
        window.filterState.type = category;

        // Close panel via the global accessible toggle function
        toggleFilterPanel(false);
        
        // Update active state in UI
        document.querySelectorAll('[data-category-name]').forEach(b => {
            b.classList.remove('bg-earth-ochre/10', 'text-earth-ochre');
        });
        btn.classList.add('bg-earth-ochre/10', 'text-earth-ochre');

        // Trigger HTMX
        const searchInput = document.getElementById('search-nav') || document.getElementById('search');
        if (searchInput) {
            if (window.htmx) {
                window.htmx.trigger(searchInput, 'search');
            } else {
                searchInput.dispatchEvent(new Event('search', { bubbles: true }));
            }
        } else {
            const city = document.getElementById('filter-city')?.value || '';
            const radius = document.getElementById('filter-radius')?.value || '25';
            const url = `/listings/fragment?type=${encodeURIComponent(category)}&city=${encodeURIComponent(city)}&radius=${encodeURIComponent(radius)}`;
            if (window.htmx) {
                window.htmx.ajax('GET', url, {
                    target: '#listings-container',
                    indicator: '#listings-loading'
                });
            }
        }
    });

    // Handle City and Radius updates to global state
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-radius-value]');
        if (!btn) return;

        const radius = btn.getAttribute('data-radius-value') || '25';
        
        // Update global state
        window.filterState.radius = radius;

        // Close panel via the global accessible toggle function
        toggleFilterPanel(false);
        
        // Update active state in UI
        document.querySelectorAll('[data-radius-value]').forEach(b => {
            b.classList.remove('bg-earth-ochre/10', 'text-earth-ochre');
        });
        btn.classList.add('bg-earth-ochre/10', 'text-earth-ochre');

        // Trigger HTMX
        const searchInput = document.getElementById('search-nav') || document.getElementById('search');
        if (searchInput) {
            if (window.htmx) {
                window.htmx.trigger(searchInput, 'search');
            } else {
                searchInput.dispatchEvent(new Event('search', { bubbles: true }));
            }
        } else {
            const type = window.filterState.type || 'Food';
            const city = document.getElementById('filter-city')?.value || '';
            const url = `/listings/fragment?type=${encodeURIComponent(type)}&city=${encodeURIComponent(city)}&radius=${encodeURIComponent(radius)}`;
            if (window.htmx) {
                window.htmx.ajax('GET', url, {
                    target: '#listings-container',
                    indicator: '#listings-loading'
                });
            }
        }
    });

    document.addEventListener('input', (e) => {
        if (e.target.id === 'filter-city') {
            window.filterState.city = e.target.value;
        }
    });

    // Listen for "VIEW ALL LISTINGS" click to reset filters & categories
    document.addEventListener('click', (e) => {
        const resetBtn = e.target.closest('[data-testid="view-all-listings-link"]');
        if (!resetBtn) return;

        // Reset global filter state
        window.filterState = {
            type: 'All',
            city: '',
            radius: '25'
        };

        // Clear UI inputs
        const searchInput = document.getElementById('search-nav') || document.getElementById('search');
        if (searchInput) searchInput.value = '';
        
        const cityInput = document.getElementById('filter-city');
        if (cityInput) cityInput.value = '';

        // Reset category button UI classes
        document.querySelectorAll('[data-category-name]').forEach(b => {
            b.classList.remove('bg-earth-ochre/10', 'text-earth-ochre');
            if (b.getAttribute('data-category-name') === 'All') {
                b.classList.add('bg-earth-ochre/10', 'text-earth-ochre');
            }
        });

        // Reset radius button UI classes
        document.querySelectorAll('[data-radius-value]').forEach(b => {
            b.classList.remove('bg-earth-ochre/10', 'text-earth-ochre');
            if (b.getAttribute('data-radius-value') === '25') {
                b.classList.add('bg-earth-ochre/10', 'text-earth-ochre');
            }
        });
    });

    // Inject filter state into all HTMX requests to /listings/fragment
    document.body.addEventListener('htmx:configRequest', (evt) => {
        if (evt.detail.path === '/listings/fragment') {
            const state = window.filterState || {};
            if (state.type) evt.detail.parameters['type'] = state.type;
            if (state.city) evt.detail.parameters['city'] = state.city;
            if (state.radius) evt.detail.parameters['radius'] = state.radius;
        }
    });
}

