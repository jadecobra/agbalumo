// agbalumo Main Application Entry Point
window.AGBALUMO_AGENT_VERSION = "1.0.0-agent-native";

document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

function initApp() {
    // Utility UI logic
    if (typeof setupMobileBottomNav === 'function') setupMobileBottomNav();
    if (typeof setupDynamicBackgrounds === 'function') setupDynamicBackgrounds();
    
    // Component logic
    if (typeof setupFeaturedCarousel === 'function') setupFeaturedCarousel();
    if (typeof setupCustomDropdowns === 'function') setupCustomDropdowns();
    
    // Interaction/Action logic
    if (typeof setupModalClosing === 'function') setupModalClosing();
    if (typeof setupModalActions === 'function') setupModalActions();
    if (typeof setupAuthActions === 'function') setupAuthActions();
    if (typeof setupFilterButtons === 'function') setupFilterButtons();
    if (typeof setupFilterToggle === 'function') setupFilterToggle();
    setupTryButtons();
    setupSearchPills();
    initInfiniteScroll();
}

function setupSearchPills() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-search-value]');
        if (btn) {
            const searchInput = document.getElementById('search-nav') || document.getElementById('search');
            if (searchInput) {
                searchInput.value = btn.dataset.searchValue;
                searchInput.dispatchEvent(new Event('search', { bubbles: true }));
            }
        }
    });
}


// Global HTMX listener for elements that need re-init
document.body.addEventListener('htmx:afterSwap', (evt) => {
    // Specifically re-init dropdowns and dynamic backgrounds on swapped content
    if (typeof initCustomDropdownsActiveState === 'function') {
        initCustomDropdownsActiveState(evt.detail.elt);
    }
    // Re-init infinite scroll observer if the listings container or sentinel itself was swapped
    if (evt.detail.target.id === 'listings-container' || evt.detail.target.id === 'infinite-scroll-sentinel') {
        initInfiniteScroll();
    }
});

function setupTryButtons() {
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-try-query]');
        if (btn) {
            const query = btn.getAttribute('data-try-query');
            const searchInput = document.getElementById('search');
            if (searchInput) {
                searchInput.value = query;
            }
        }
    });
}

function initInfiniteScroll() {
    const sentinel = document.getElementById('infinite-scroll-sentinel');
    if (!sentinel) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting && window.innerWidth < 768) {
                sentinel.dispatchEvent(new CustomEvent('infinite-scroll'));
            }
        });
    }, { rootMargin: '100px' });

    observer.observe(sentinel);
}

