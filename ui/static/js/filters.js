/**
 * Filters Component Logic
 * Robust discovery filtering for the agbalumo platform.
 * CSP-compliant: Standard DOM listeners, avoid inline Alpine.
 */

// Global state for HTMX hx-vals
if (!window.filterState) {
    window.filterState = {
        type: 'Food',
        city: ''
    };
}

function setupFilterToggle() {
    if (window._filterToggleInitialized) return;
    window._filterToggleInitialized = true;

    const togglePanel = (show) => {
        const panel = document.getElementById('filter-dropdown-panel');
        if (!panel) return;

        if (show === undefined) {
            show = panel.classList.contains('hidden');
        }
        
        if (show) {
            panel.classList.remove('hidden');
            panel.setAttribute('aria-expanded', 'true');
            // Focus search input on open
            const searchInput = document.getElementById('search');
            if (searchInput) setTimeout(() => searchInput.focus(), 100);
        } else {
            panel.classList.add('hidden');
            panel.setAttribute('aria-expanded', 'false');
        }
    };

    document.addEventListener('click', (e) => {
        const toggle = e.target.closest('[data-testid^="ag-home-filters-toggle"]');
        if (toggle) {
            e.stopPropagation();
            togglePanel();
            return;
        }

        const panel = document.getElementById('filter-dropdown-panel');
        if (panel && !panel.contains(e.target) && !panel.classList.contains('hidden')) {
            togglePanel(false);
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') togglePanel(false);
    });
}

function setupFilterButtons() {
    if (window._filterButtonsInitialized) return;
    window._filterButtonsInitialized = true;

    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-category-name]');
        if (!btn) return;

        const category = btn.getAttribute('data-category-name') || '';
        const panel = document.getElementById('filter-dropdown-panel');
        
        // Update global state
        window.filterState.type = category;

        // Close panel
        if (panel) panel.classList.add('hidden');
        
        // Update active state in UI
        document.querySelectorAll('[data-category-name]').forEach(b => {
            b.classList.remove('bg-earth-ochre/10', 'text-earth-ochre');
        });
        btn.classList.add('bg-earth-ochre/10', 'text-earth-ochre');

        // Trigger HTMX
        const searchInput = document.getElementById('search');
        if (searchInput) {
            searchInput.dispatchEvent(new Event('search', { bubbles: true }));
        } else {
            const url = category ? `/listings/fragment?type=${encodeURIComponent(category)}` : '/listings/fragment';
            if (window.htmx) {
                window.htmx.ajax('GET', url, {
                    target: '#listings-container',
                    indicator: '#listings-loading'
                });
            }
        }
    });
}
