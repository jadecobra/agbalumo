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
                } catch (e) {
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
    const csvUploadDropzone = document.getElementById('csvUploadDropzone');
    const csvUploadIcon = document.getElementById('csvUploadIcon');
    const csvUploadText = document.getElementById('csvUploadText');

    if (csvFileInput && csvUploadBtn) {
        csvFileInput.addEventListener('change', function () {
            const hasFile = this.files && this.files.length > 0;
            csvUploadBtn.disabled = !hasFile;
            
            if (hasFile && csvUploadDropzone && csvUploadIcon && csvUploadText) {
                const fileName = this.files[0].name;
                
                // Update dropzone styling
                csvUploadDropzone.classList.remove('border-dashed', 'border-white/10', 'bg-white/2');
                csvUploadDropzone.classList.add('border-solid', 'border-earth-ochre', 'bg-earth-ochre/10');
                
                // Update icon
                csvUploadIcon.textContent = 'description';
                csvUploadIcon.classList.remove('text-white/20');
                csvUploadIcon.classList.add('text-earth-ochre');
                
                // Update text securely without innerHTML
                csvUploadText.textContent = '';
                
                const nameSpan = document.createElement('span');
                nameSpan.className = 'text-white text-sm tracking-normal mb-1 block normal-case font-normal';
                nameSpan.textContent = fileName;
                
                const clickSpan = document.createElement('span');
                clickSpan.className = 'text-[9px] text-earth-ochre font-bold uppercase tracking-[0.2em]';
                clickSpan.textContent = 'Click to replace';

                csvUploadText.appendChild(nameSpan);
                csvUploadText.appendChild(clickSpan);
            } else if (csvUploadDropzone && csvUploadIcon && csvUploadText) {
                // Revert to initial state
                csvUploadDropzone.classList.add('border-dashed', 'border-white/10', 'bg-white/2');
                csvUploadDropzone.classList.remove('border-solid', 'border-earth-ochre', 'bg-earth-ochre/10');
                
                csvUploadIcon.textContent = 'add';
                csvUploadIcon.classList.add('text-white/20');
                csvUploadIcon.classList.remove('text-earth-ochre');
                
                csvUploadText.textContent = 'Choose CSV File';
                csvUploadText.className = 'text-[10px] font-bold uppercase tracking-[0.2em] text-white/60';
            }
        });
    }

    // Modal Handling
    document.addEventListener('click', function (e) {
        // Open Modal
        let openTrigger = e.target.closest('[data-action="open-modal"]');
        if (openTrigger) {
            e.preventDefault();
            const targetId = openTrigger.dataset.target;
            const targetModal = document.getElementById(targetId);
            if (targetModal) {
                targetModal.classList.remove('hidden');

                // If opening charts modal, we might need to force a reflow/update 
                // for Chart.js if they were initialized while hidden.
                // Chart.instances are global if stored, but 'responsive: true' usually handles it.
            }
        }

        // Close Modal
        let closeTrigger = e.target.closest('[data-action="close-modal"]');
        if (closeTrigger) {
            e.preventDefault();
            const targetId = closeTrigger.dataset.target;
            const targetModal = document.getElementById(targetId);
            if (targetModal) {
                targetModal.classList.add('hidden');
            }
        }
    });
}

document.addEventListener('DOMContentLoaded', initAdminDashboard);
