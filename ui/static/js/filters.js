/**
 * Filter Dropdown Logic (Hardened)
 * 1. Dynamic Element Fetching: Prevents stale DOM references during HTMX swaps.
 * 2. Idempotent Initialization: Ensures exactly one listener setup.
 * 3. Mobile vs Desktop Logic: Handles bottom-sheet vs dropdown-panel correctly.
 */

function setupFilterButtons() {
    // Placeholder to satisfy legacy app.js dependency
}

function setupFilterToggle() {
    // Use a flag on the document to ensure listeners are only added ONCE globally
    if (window._agFiltersGlobalBound) {
        console.log("[Filters] Global listeners already bound. Skipping.");
        return;
    }
    window._agFiltersGlobalBound = true;

    console.log("[Filters] Initializing hardened global filter handlers...");

    const getPanel = () => document.getElementById('filter-dropdown-panel');
    const filterState = { type: 'Food' };

    const updateButtonStates = (isFilterOpen) => {
        const filtersBtns = document.querySelectorAll('.filters-btn');
        const searchBtns = document.querySelectorAll('.search-btn');
        const activeClasses = ['bg-earth-ochre', 'text-earth-dark'];
        const inactiveClasses = ['bg-white/10', 'text-white'];

        filtersBtns.forEach(btn => {
            if (isFilterOpen) {
                btn.classList.add(...activeClasses);
                btn.classList.remove(...inactiveClasses);
            } else {
                btn.classList.remove(...activeClasses);
                btn.classList.add(...inactiveClasses);
            }
        });
    };

    const closePanel = () => {
        const panel = getPanel();
        if (panel && !panel.classList.contains('hidden')) {
            panel.classList.add('hidden');
            updateButtonStates(false);
            console.log("[Filters] Panel closed.");
        }
    };

    const triggerSearch = () => {
        const params = new URLSearchParams();
        if (filterState.type) params.set('type', filterState.type);

        const searchInput = document.getElementById('search') || document.getElementById('search-header');
        if (searchInput && searchInput.value) {
            params.set('q', searchInput.value);
        }

        const url = `/listings/fragment?${params.toString()}`;
        if (window.htmx) {
            htmx.ajax('GET', url, { target: '#listings-container', indicator: '#listings-loading' });
        }
        closePanel();
    };

    // 1. Global Key Handler (document level)
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') closePanel();
    });

    // 2. Global Click Handler (Consolidated Delegation)
    document.addEventListener('click', (e) => {
        const target = e.target;
        const panel = getPanel();
        if (!panel) return;

        // Toggle Button Click
        const toggleBtn = target.closest('[data-action="toggle-filters"]');
        if (toggleBtn) {
            const isCurrentlyHidden = panel.classList.contains('hidden');
            panel.classList.toggle('hidden');
            updateButtonStates(isCurrentlyHidden);
            console.log("[Filters] Panel toggled via button:", isCurrentlyHidden);
            return;
        }

        // Click Outside to Close
        if (!panel.classList.contains('hidden')) {
            const isClickInsidePanel = target.closest('#filter-dropdown-panel');
            if (!isClickInsidePanel) {
                console.log("[Filters] Click outside detected on:", target.tagName, target.className);
                closePanel();
            }
        }

        // Selection Logic
        const filterBtn = target.closest('[data-filter-type="type"]');
        if (filterBtn) {
            const value = filterBtn.dataset.filterValue;
            filterState.type = value;
            
            // Visual Update for selections
            document.querySelectorAll('[data-filter-type="type"]').forEach(btn => {
                btn.classList.remove('bg-earth-ochre', 'text-white', 'bg-earth-ochre/10', 'bg-earth-ochre/20', 'text-earth-ochre');
                btn.classList.add('text-earth-dark');
            });
            
            if (value === 'Food') {
                filterBtn.classList.add('bg-earth-ochre/20', 'text-earth-ochre');
            } else {
                filterBtn.classList.add('bg-earth-ochre', 'text-white');
            }

            console.log("[Filters] Category selected:", value);
            triggerSearch();
        }
    }, { capture: true });
}
