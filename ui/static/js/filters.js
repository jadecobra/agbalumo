// Combined filter state - keys match backend query parameters
window.filterState = {
    type: 'Food', // Default for Ada
    city: ''
};

function setupFilterToggle() {
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

    // Initialize state from existing active classes in DOM on load
    document.querySelectorAll('[data-filter-type]').forEach(btn => {
        if (btn.classList.contains('bg-earth-ochre')) {
            filterState[btn.dataset.filterType] = btn.dataset.filterValue;
        }
    });

    document.addEventListener('click', (e) => {
        const toggleBtn = e.target.closest('[data-action="toggle-filters"]');
        if (toggleBtn) {
            const panel = document.getElementById('filter-dropdown-panel');
            if (panel) {
                const isOpen = !panel.classList.contains('hidden');
                panel.classList.toggle('hidden');
                updateButtonStates(!isOpen);
            }
            return;
        }

        const chip = e.target.closest('[data-filter-type]');
        if (chip) {
            const type = chip.dataset.filterType;
            const value = chip.dataset.filterValue;

            // Update state
            filterState[type] = value;

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

            // Trigger HTMX request
            const params = new URLSearchParams();
            if (filterState.type) params.set('type', filterState.type);
            if (filterState.city) params.set('city', filterState.city);
            
            // Gather search query from ANY search input on the page
            const searchInput = document.getElementById('search') || document.getElementById('search-header');
            if (searchInput && searchInput.value) {
                params.set('q', searchInput.value);
            }

            const url = `/listings/fragment?${params.toString()}`;
            htmx.ajax('GET', url, { target: '#listings-container', indicator: '#listings-loading' });
            
            // Auto-close panel on selection (optional, following existing pattern)
            const panel = document.getElementById('filter-dropdown-panel');
            if (panel) panel.classList.add('hidden');
            updateButtonStates(false);
        }
    });
}

// Global initialization
document.addEventListener('DOMContentLoaded', () => {
    setupFilterToggle();
});
