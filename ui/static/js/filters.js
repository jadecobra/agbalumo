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
            const items = document.querySelectorAll('.location-item');
            const groups = document.querySelectorAll('.location-group-content');
            const toggles = document.querySelectorAll('.location-group-toggle');

            groups.forEach(group => {
                const groupItems = group.querySelectorAll('.location-item');
                let hasMatch = false;

                groupItems.forEach(item => {
                    const name = (item.dataset.cityName || '').toLowerCase();
                    const isMatch = name.includes(val);
                    item.classList.toggle('hidden', !isMatch);
                    if (isMatch) hasMatch = true;
                });

                // Auto-expand if has match and search is not empty
                if (val.length > 0 && hasMatch) {
                    group.classList.remove('max-h-0', 'opacity-0');
                    group.classList.add('max-h-screen', 'opacity-100');
                    const toggle = document.querySelector(`.location-group-toggle[data-group="${group.dataset.groupName}"]`);
                    if (toggle) {
                        const icon = toggle.querySelector('[data-toggle-icon]');
                        if (icon) icon.classList.add('rotate-180');
                    }
                } else if (val.length === 0) {
                    // Reset to collapsed if cleared (optional, but cleaner)
                    group.classList.add('max-h-0', 'opacity-0');
                    group.classList.remove('max-h-screen', 'opacity-100');
                    const toggle = document.querySelector(`.location-group-toggle[data-group="${group.dataset.groupName}"]`);
                    if (toggle) {
                        const icon = toggle.querySelector('[data-toggle-icon]');
                        if (icon) icon.classList.remove('rotate-180');
                    }
                }

                // Hide the toggle (header) if no matches in group and searching
                const toggle = document.querySelector(`.location-group-toggle[data-group="${group.dataset.groupName}"]`);
                if (toggle) {
                    toggle.classList.toggle('hidden', val.length > 0 && !hasMatch);
                }
            });
        });
    }

    document.addEventListener('click', (e) => {
        const target = e.target;

        // Location Group Accordion Toggle
        const groupToggle = target.closest('.location-group-toggle');
        if (groupToggle) {
            const groupName = groupToggle.dataset.group;
            const content = document.querySelector(`.location-group-content[data-group-name="${groupName}"]`);
            const icon = groupToggle.querySelector('[data-toggle-icon]');
            
            if (content) {
                const isExpanded = content.classList.contains('max-h-screen');
                // Toggle classes
                content.classList.toggle('max-h-0', isExpanded);
                content.classList.toggle('opacity-0', isExpanded);
                content.classList.toggle('max-h-screen', !isExpanded);
                content.classList.toggle('opacity-100', !isExpanded);
                
                if (icon) {
                    icon.classList.toggle('rotate-180', !isExpanded);
                }
            }
            return;
        }

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
