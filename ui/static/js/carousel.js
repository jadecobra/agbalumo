function setupFeaturedCarousel() {
    const carousel = document.getElementById('featured-carousel');
    if (!carousel) return;

    const slides = carousel.querySelectorAll('.carousel-slide');
    const dots = carousel.querySelectorAll('.carousel-dot');
    const prevBtn = carousel.querySelector('.carousel-prev');
    const nextBtn = carousel.querySelector('.carousel-next');

    if (slides.length <= 1) return;

    let currentIndex = 0;
    let interval = null;
    let isPaused = false;

    function goToSlide(index) {
        slides.forEach((slide, i) => {
            slide.classList.toggle('opacity-100', i === index);
            slide.classList.toggle('z-10', i === index);
            slide.classList.toggle('opacity-0', i !== index);
            slide.classList.toggle('z-0', i !== index);
        });

        dots.forEach((dot, i) => {
            dot.classList.toggle('bg-white', i === index);
            dot.classList.toggle('w-6', i === index);
            dot.classList.toggle('md:w-8', i === index);
            dot.classList.toggle('bg-white/40', i !== index);
            dot.classList.toggle('w-2', i !== index);
            dot.classList.toggle('md:w-3', i !== index);
        });

        currentIndex = index;
    }

    function nextSlide() {
        goToSlide((currentIndex + 1) % slides.length);
    }

    function prevSlide() {
        goToSlide((currentIndex - 1 + slides.length) % slides.length);
    }

    function startAutoplay() {
        if (interval) clearInterval(interval);
        const delay = 5000 + Math.random() * 3000;
        interval = setInterval(() => {
            if (!isPaused) nextSlide();
        }, delay);
    }

    if (prevBtn) prevBtn.onclick = (e) => { e.stopPropagation(); prevSlide(); startAutoplay(); };
    if (nextBtn) nextBtn.onclick = (e) => { e.stopPropagation(); nextSlide(); startAutoplay(); };

    dots.forEach((dot, i) => {
        dot.onclick = (e) => { e.stopPropagation(); goToSlide(i); startAutoplay(); };
    });

    carousel.onmouseenter = () => isPaused = true;
    carousel.onmouseleave = () => isPaused = false;
    carousel.onfocusin = () => isPaused = true;
    carousel.onfocusout = () => isPaused = false;

    if (carousel.dataset.autoplay === 'true') startAutoplay();
}
