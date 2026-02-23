function initAdminDashboard() {
    const listingChartEl = document.getElementById('listingChart');
    const userChartEl = document.getElementById('userChart');
    
    if (listingChartEl || userChartEl) {
        if (typeof Chart === 'undefined') {
            console.warn('Chart.js not loaded');
        } else {
            function parseDataFromAttr(el) {
                if (!el || !el.dataset.chartData) return { labels: [], values: [] };
                try {
                    const data = JSON.parse(el.dataset.chartData);
                    return {
                        labels: data.map(d => d.Date),
                        values: data.map(d => d.Count)
                    };
                } catch(e) {
                    return { labels: [], values: [] };
                }
            }

            const listings = parseDataFromAttr(listingChartEl);
            const users = parseDataFromAttr(userChartEl);

            if (listingChartEl) {
                new Chart(listingChartEl, {
                    type: 'line',
                    data: {
                        labels: listings.labels,
                        datasets: [{
                            label: 'New Listings',
                            data: listings.values,
                            borderColor: 'rgb(234, 88, 12)',
                            backgroundColor: 'rgba(234, 88, 12, 0.1)',
                            fill: true,
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: { legend: { display: false } },
                        scales: { y: { beginAtZero: true } }
                    }
                });
            }

            if (userChartEl) {
                new Chart(userChartEl, {
                    type: 'line',
                    data: {
                        labels: users.labels,
                        datasets: [{
                            label: 'New Users',
                            data: users.values,
                            borderColor: 'rgb(34, 197, 94)',
                            backgroundColor: 'rgba(34, 197, 94, 0.1)',
                            fill: true,
                        }]
                    },
                    options: {
                        responsive: true,
                        plugins: { legend: { display: false } },
                        scales: { y: { beginAtZero: true } }
                    }
                });
            }
        }
    }

    const csvFileInput = document.getElementById('csvFileInput');
    const csvUploadBtn = document.getElementById('csvUploadBtn');
    if (csvFileInput && csvUploadBtn) {
        csvFileInput.addEventListener('change', function () {
            csvUploadBtn.disabled = !(this.files && this.files.length > 0);
        });
    }
}

document.addEventListener('DOMContentLoaded', initAdminDashboard);
