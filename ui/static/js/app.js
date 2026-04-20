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
}

// Global HTMX listener for elements that need re-init
document.body.addEventListener('htmx:afterSwap', (evt) => {
    // Specifically re-init dropdowns and dynamic backgrounds on swapped content
    if (typeof initCustomDropdownsActiveState === 'function') {
        initCustomDropdownsActiveState(evt.detail.elt);
    }
});
