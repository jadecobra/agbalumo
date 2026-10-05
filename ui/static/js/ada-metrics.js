(function() {
    const ADA_SESSION_START = 'ada_session_start';
    const DISCOVERY_EVENT = 'discovery_success';

    if (!sessionStorage.getItem(ADA_SESSION_START)) {
        sessionStorage.setItem(ADA_SESSION_START, Date.now());
    }

    const gaMeta = document.querySelector('meta[name="ga-measurement-id"]');
    const gaMeasurementId = gaMeta ? gaMeta.getAttribute('content') : null;

    if (gaMeasurementId && !window.gtag) {
        window.dataLayer = window.dataLayer || [];
        window.gtag = function() {
            window.dataLayer.push(arguments);
        };
        window.gtag('js', new Date());
        window.gtag('config', gaMeasurementId, {
            send_page_view: true
        });

        const script = document.createElement('script');
        script.async = true;
        script.src = `https://www.googletagmanager.com/gtag/js?id=${encodeURIComponent(gaMeasurementId)}`;
        document.head.appendChild(script);
    }

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

    document.addEventListener('click', (e) => {
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
                
                if (!sessionStorage.getItem('ada_discovered')) {
                    sessionStorage.setItem('ada_discovered', 'true');
                    sendMetric('first_discovery_success', duration, metadata);
                }
            }
        }
    });

    let currentPath = window.location.pathname;
    document.body.addEventListener('htmx:afterSwap', () => {
        if (window.location.pathname !== currentPath) {
            currentPath = window.location.pathname;
            if (typeof window.gtag === 'function' && gaMeasurementId) {
                window.gtag('event', 'page_view', {
                    page_path: currentPath,
                    page_location: window.location.href,
                    page_title: document.title
                });
            }
        }
    });

    async function sendMetric(event, value, metadata) {
        if (typeof window.gtag === 'function') {
            try {
                window.gtag('event', event, {
                    value: value,
                    ...(metadata || {})
                });
            } catch (_) {}
        }

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
        } catch (err) {}
    }
})();
