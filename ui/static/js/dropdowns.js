function initCustomDropdownsActiveState(root = document) {
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
        const toggleBtn = e.target.closest('[data-dropdown-toggle]');
        if (toggleBtn) {
            const container = toggleBtn.closest('.custom-dropdown');
            if (!container) return;
            const menu = container.querySelector('.dropdown-menu');
            if (menu) {
                const isHidden = menu.classList.contains('hidden');
                
                document.querySelectorAll('.custom-dropdown .dropdown-menu').forEach(m => {
                    if (m !== menu) {
                        m.classList.add('hidden');
                        const otherBtn = m.closest('.custom-dropdown').querySelector('[data-dropdown-toggle]');
                        if (otherBtn) otherBtn.setAttribute('aria-expanded', 'false');
                    }
                });
                
                menu.classList.toggle('hidden');
                toggleBtn.setAttribute('aria-expanded', isHidden ? 'true' : 'false');
                
                document.querySelectorAll('.custom-dropdown').forEach(dropdown => {
                    dropdown.style.zIndex = dropdown === container ? '50' : '10';
                });
            }
            return;
        }

        const optionBtn = e.target.closest('[data-dropdown-value]');
        if (optionBtn) {
            const container = optionBtn.closest('.custom-dropdown');
            if (!container) return;

            const value = optionBtn.dataset.dropdownValue;
            let label = optionBtn.textContent.trim();
            const input = container.querySelector('input[type="hidden"]');
            const display = container.querySelector('.dropdown-display');
            const menu = container.querySelector('.dropdown-menu');
            const toggle = container.querySelector('[data-dropdown-toggle]');

            if (input && input.value !== value) {
                input.value = value;
                input.dispatchEvent(new Event('change', { bubbles: true }));
            }
            if (display) display.textContent = label;

            container.querySelectorAll('[data-dropdown-value]').forEach(btn => {
                btn.classList.remove('bg-earth-ochre', 'text-white');
                btn.classList.add('bg-earth-dark/5', 'text-earth-dark');
                btn.setAttribute('aria-selected', 'false');
            });
            optionBtn.classList.remove('bg-earth-dark/5', 'text-earth-dark');
            optionBtn.classList.add('bg-earth-ochre', 'text-white');
            optionBtn.setAttribute('aria-selected', 'true');

            if (menu) {
                menu.classList.add('hidden');
                if (toggle) toggle.setAttribute('aria-expanded', 'false');
                container.style.zIndex = '10';
            }
            return;
        }

        if (!e.target.closest('.custom-dropdown')) {
            document.querySelectorAll('.custom-dropdown .dropdown-menu').forEach(m => {
                m.classList.add('hidden');
                const container = m.closest('.custom-dropdown');
                if (container) {
                    container.style.zIndex = '10';
                    const toggle = container.querySelector('[data-dropdown-toggle]');
                    if (toggle) toggle.setAttribute('aria-expanded', 'false');
                }
            });
        }
    });
}
