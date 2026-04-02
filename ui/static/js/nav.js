function setupMobileBottomNav() {
    const nav = document.getElementById('mobile-bottom-nav');
    if (!nav) return;

    const scrollContainer = document.querySelector('main');
    if (!scrollContainer) return;

    let lastScrollY = 0;
    let scrollTimeout;

    scrollContainer.addEventListener('scroll', () => {
        const currentScrollY = scrollContainer.scrollTop;
        const isScrollingDown = currentScrollY > lastScrollY && currentScrollY > 60;

        if (isScrollingDown) {
            nav.classList.add('nav-hidden');
        } else {
            nav.classList.remove('nav-hidden');
        }

        lastScrollY = currentScrollY;

        clearTimeout(scrollTimeout);
        scrollTimeout = setTimeout(() => {
            nav.classList.remove('nav-hidden');
        }, 1500);
    }, { passive: true });
}
