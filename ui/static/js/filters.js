// Combined filter state - keys match backend query parameters
window.filterState = {
    type: 'Food', // Default for Ada
    city: ''
};

function setupFilterToggle() {
    console.log('[Filters] Initializing filter logic...');
    const activeClasses = ['bg-earth-ochre', 'text-white'];
    const inactiveClasses = ['bg-earth-dark/5', 'text-earth-dark'];

    const updateButtonStates = (isFilterOpen) => {
        const filtersBtns = document.querySelectorAll('.filters-btn');
        const searchBtns = document.querySelectorAll('.search-btn');
        const btnActive = ['bg-earth-ochre', 'text-earth-dark'];
        const btnInactive = ['bg-white/10', 'text-white'];

        if (isFilterOpen) {
            filtersBtns.forEach(btn => {
                btn.classList.add(...btnActive);
                btn.classList.remove(...btnInactive);
            });
            searchBtns.forEach(btn => {
                btn.classList.remove(...btnActive);
                btn.classList.add(...btnInactive);
            });
        } else {
            filtersBtns.forEach(btn => {
                btn.classList.remove(...btnActive);
                btn.classList.add(...btnInactive);
            });
            searchBtns.forEach(btn => {
                btn.classList.add(...btnActive);
                btn.classList.remove(...btnInactive);
            });
        }
    };

    const triggerSearch = () => {
        const params = new URLSearchParams();
        if (filterState.type) params.set('type', filterState.type);
        if (filterState.city) params.set('city', filterState.city);

        // Gather search query from ANY search input on the page
        const searchInput = document.getElementById('search') || document.getElementById('search-header');
        if (searchInput && searchInput.value) {
            params.set('q', searchInput.value);
        }

        console.log('[Filters] Triggering search with params:', params.toString());
        const url = `/listings/fragment?${params.toString()}`;
        htmx.ajax('GET', url, { target: '#listings-container', indicator: '#listings-loading' });

        // Auto-close panel on selection
        const panel = document.getElementById('filter-dropdown-panel');
        if (panel) panel.classList.add('hidden');
        updateButtonStates(false);
    };

    // Initialize state from existing active classes in DOM on load
    document.querySelectorAll('[data-filter-type]').forEach(btn => {
        if (btn.classList.contains('bg-earth-ochre')) {
            filterState[btn.dataset.filterType] = btn.dataset.filterValue;
        }
    });

    // Handle Location Search (Search-as-you-type)
    const locSearch = document.getElementById('location-search');
    if (locSearch) {
        console.log('[Filters] Location search input found, attaching listener');
        locSearch.addEventListener('input', (e) => {
            const val = e.target.value.trim().toLowerCase();
            console.log('[Filters] Location filtering for:', val);
            const items = document.querySelectorAll('.location-item');
            const headers = document.querySelectorAll('.location-header');

            // Filter items
            items.forEach(item => {
                const name = (item.dataset.cityName || '').toLowerCase();
                if (name.includes(val)) {
                    item.classList.remove('hidden');
                } else {
                    item.classList.add('hidden');
                }
            });

            // Filter headers based on item visibility
            headers.forEach(header => {
                let hasVisible = false;
                let sibling = header.nextElementSibling;
                while (sibling && !sibling.classList.contains('location-header')) {
                    if (sibling.classList.contains('location-item') && !sibling.classList.contains('hidden')) {
                        hasVisible = true;
                        break;
                    }
                    sibling = sibling.nextElementSibling;
                }
                header.classList.toggle('hidden', !hasVisible);
            });
        });
    }

    document.addEventListener('click', (e) => {
        const target = e.target;

        // Toggle Filter Panel
        const toggleBtn = target.closest('[data-action="toggle-filters"]');
        if (toggleBtn) {
            const panel = document.getElementById('filter-dropdown-panel');
            if (panel) {
                const isOpen = !panel.classList.contains('hidden');
                panel.classList.toggle('hidden');
                updateButtonStates(!isOpen);
            }
            return;
        }

        // Search Action
        const searchBtn = target.closest('.search-btn');
        if (searchBtn) {
            triggerSearch();
            return;
        }

        // Filter Chip Selection
        const chip = target.closest('[data-filter-type]');
        if (chip) {
            const type = chip.dataset.filterType;
            const value = chip.dataset.filterValue;

            console.log(`[Filters] Chip selected: ${type}=${value}`);

            // Update global state
            filterState[type] = value;

            // Update main search input if appropriate
            if (type === 'city') {
                const mainSearch = document.getElementById('search');
                if (mainSearch) {
                    mainSearch.value = value; // Sync city name to search bar
                }
            }

            // Update UI for the group
            const group = chip.closest('.flex-col');
            if (group) {
                group.querySelectorAll('[data-filter-type]').forEach(btn => {
                    btn.classList.remove(...activeClasses);
                    btn.classList.add(...inactiveClasses);
                });
                chip.classList.add(...activeClasses);
                chip.classList.remove(...inactiveClasses);
            }

            triggerSearch();
        }
    });
}

// Global initialization now managed centrally by app.js initApp()
