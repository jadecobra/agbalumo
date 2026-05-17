(function() {
    const ADA_SESSION_START = 'ada_session_start';
    const DISCOVERY_EVENT = 'discovery_success';

    // Initialize session start time if not present
    if (!sessionStorage.getItem(ADA_SESSION_START)) {
        sessionStorage.setItem(ADA_SESSION_START, Date.now());
    }

    // Capture exact page load time from navigationStart
    const navStart = (window.performance && window.performance.timing) ? window.performance.timing.navigationStart : Date.now();
    const pageLoadDurationMs = Date.now() - navStart;

    // Capture contact clicks
    document.addEventListener('click', (e) => {
        // We look for any link or button with data-ada-discovery
        const contactLink = e.target.closest('[data-ada-discovery]');
        if (contactLink) {
            const clickTime = Date.now();
            const durationFromStartSec = (clickTime - navStart) / 1000;
            
            // Standard discovery click success signal
            sendMetric(DISCOVERY_EVENT, durationFromStartSec, {
                type: contactLink.dataset.adaDiscovery,
                path: window.location.pathname,
                page_load_ms: pageLoadDurationMs
            });
            
            // For the 60s goal, we care about the FIRST one in the session.
            if (!sessionStorage.getItem('ada_discovered')) {
                sessionStorage.setItem('ada_discovered', 'true');
                sendMetric('first_discovery_success', durationFromStartSec, {
                    type: contactLink.dataset.adaDiscovery,
                    path: window.location.pathname,
                    page_load_ms: pageLoadDurationMs
                });
            }
        }
    });

    async function sendMetric(event, value, metadata) {
        try {
            const token = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
            await fetch('/api/metrics', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': token
                },
                body: JSON.stringify({ event, value, metadata })
            });
        } catch (err) {
            // Silently fail to not disturb user
        }
    }
})();
