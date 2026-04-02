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
                if (icon && !icon.classList.contains('text-text-sub')) {
                    icon.classList.add('text-text-sub');
                }
            });

            btn.className = activeState;
            const icon = btn.querySelector('.material-symbols-outlined');
            if (icon) icon.classList.remove('text-text-sub');
        }
    }, true);
}

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

                if (window.innerWidth < 768) {
                    const bottomNav = document.getElementById('mobile-bottom-nav');
                    let overlay = document.getElementById('mobile-filter-overlay');
                    if (!overlay) {
                        overlay = document.createElement('div');
                        overlay.id = 'mobile-filter-overlay';
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
                if (window.innerWidth < 768) document.body.style.overflow = '';
            }
            if (searchBtn) {
                updateButtonStates(false);
                return;
            }
        }

        const chip = e.target.closest('[data-filter-type]');
        if (chip) {
            const filterType = chip.dataset.filterType;
            const filterValue = chip.dataset.filterValue;
            const container = chip.parentElement;

            if (container) {
                container.querySelectorAll('button').forEach(b => {
                    b.classList.remove('bg-earth-ochre', 'text-white');
                    b.classList.add('bg-earth-dark/5', 'text-earth-dark');
                });
                chip.classList.add('bg-earth-ochre', 'text-white');
                chip.classList.remove('bg-earth-dark/5', 'text-earth-dark');
            }

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
