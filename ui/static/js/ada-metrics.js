(function() {
    const ADA_SESSION_START = 'ada_session_start';
    const DISCOVERY_EVENT = 'discovery_success';

    // Initialize session start time if not present
    if (!sessionStorage.getItem(ADA_SESSION_START)) {
        sessionStorage.setItem(ADA_SESSION_START, Date.now());
    }

    // Capture inbound campaign parameters if present
    const urlParams = new URLSearchParams(window.location.search);
    const utmSource = urlParams.get('utm_source');
    if (utmSource) {
        sessionStorage.setItem('utm_source', utmSource);
        const utmMedium = urlParams.get('utm_medium') || '';
        const utmCampaign = urlParams.get('utm_campaign') || '';
        sessionStorage.setItem('utm_medium', utmMedium);
        sessionStorage.setItem('utm_campaign', utmCampaign);

        if (!sessionStorage.getItem('ada_campaign_tracked')) {
            sessionStorage.setItem('ada_campaign_tracked', 'true');
            sendMetric('campaign_visit', 1, {
                utm_source: utmSource,
                utm_medium: utmMedium,
                utm_campaign: utmCampaign,
                landing_path: window.location.pathname
            });
        }
    }

    // Capture contact clicks
    document.addEventListener('click', (e) => {
        // We look for any link or button with data-ada-discovery
        const contactLink = e.target.closest('[data-ada-discovery]');
        if (contactLink) {
            const startTime = sessionStorage.getItem(ADA_SESSION_START);
            if (startTime) {
                const duration = (Date.now() - startTime) / 1000;
                const metadata = {
                    type: contactLink.dataset.adaDiscovery,
                    path: window.location.pathname
                };
                const source = sessionStorage.getItem('utm_source');
                if (source) {
                    metadata.utm_source = source;
                    metadata.utm_medium = sessionStorage.getItem('utm_medium') || '';
                    metadata.utm_campaign = sessionStorage.getItem('utm_campaign') || '';
                }
                
                sendMetric(DISCOVERY_EVENT, duration, metadata);
                
                // To measure "First discovery", we could clear the session start,
                // but usually we want to see if they find multiple things.
                // For the 60s goal, we care about the FIRST one.
                // Let's add a "first" flag if they haven't discovered yet.
                if (!sessionStorage.getItem('ada_discovered')) {
                    sessionStorage.setItem('ada_discovered', 'true');
                    sendMetric('first_discovery_success', duration, metadata);
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
