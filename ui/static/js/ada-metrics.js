(function() {
    const ADA_SESSION_START = 'ada_session_start';
    const DISCOVERY_EVENT = 'discovery_success';

    // Initialize session start time if not present
    if (!sessionStorage.getItem(ADA_SESSION_START)) {
        sessionStorage.setItem(ADA_SESSION_START, Date.now());
    }

    // Capture contact clicks
    document.addEventListener('click', (e) => {
        // We look for any link or button with data-ada-discovery
        const contactLink = e.target.closest('[data-ada-discovery]');
        if (contactLink) {
            const startTime = sessionStorage.getItem(ADA_SESSION_START);
            if (startTime) {
                const duration = (Date.now() - startTime) / 1000;
                
                // If it's the first discovery in this session, we mark it specially?
                // For now, let's keep it simple: every discovery click is a success signal.
                
                sendMetric(DISCOVERY_EVENT, duration, {
                    type: contactLink.dataset.adaDiscovery,
                    path: window.location.pathname
                });
                
                // To measure "First discovery", we could clear the session start,
                // but usually we want to see if they find multiple things.
                // For the 60s goal, we care about the FIRST one.
                // Let's add a "first" flag if they haven't discovered yet.
                if (!sessionStorage.getItem('ada_discovered')) {
                    sessionStorage.setItem('ada_discovered', 'true');
                    sendMetric('first_discovery_success', duration, {
                         type: contactLink.dataset.adaDiscovery
                    });
                }
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
