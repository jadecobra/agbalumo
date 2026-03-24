// Listing Creation and Form Validation logic

document.addEventListener('DOMContentLoaded', () => {
    setupListingModalDelegation();
    setupCreateImagePreviewInit();
});

// Create Image Preview Init (replaces inline script)
function setupCreateImagePreviewInit() {
    // Init on page load if create modal already exists in DOM
    if (document.getElementById('create-listing-modal') && window.initCreateImagePreview) {
        window.initCreateImagePreview();
    }

    // Also init when create modal is opened via button click
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-modal-id="create-listing-modal"]');
        if (btn) {
            setTimeout(() => {
                if (window.initCreateImagePreview) window.initCreateImagePreview();
            }, 50);
        }
    });
}

// 3. Create Listing Modal Logic
// We attach this to document because the modal is inserted dynamically or exists statically but hidden.
function setupListingModalDelegation() {
    document.addEventListener('change', (event) => {
        if (event.target.matches('#create-listing-modal select[name="type"], #create-listing-modal input[name="type"]')) {
            toggleListingFields(event.target);
        }
    });

    const typeSelect = document.querySelector('#create-listing-modal select[name="type"], #create-listing-modal input[name="type"]');
    if (typeSelect) {
        toggleListingFields(typeSelect);
    }
}

function toggleListingFields(typeSelect) {
    const modal = typeSelect.closest('dialog');
    if (!modal) return;

    const eventSection = modal.querySelector('#event-dates-section');
    const jobSection = modal.querySelector('#job-fields-section');
    const imageSection = modal.querySelector('#image-upload-section');
    const hoursSection = modal.querySelector('#hours-section');
    const locationLabel = modal.querySelector('#location-label');
    const descriptionLabel = modal.querySelector('#description-label');
    const addressInput = modal.querySelector('#create-address-input');

    const val = typeSelect.value;

    // Event Logic
    if (eventSection) {
        if (val === 'Event') {
            eventSection.classList.remove('hidden');
            eventSection.querySelectorAll('input').forEach(i => i.required = true);
        } else {
            eventSection.classList.add('hidden');
            eventSection.querySelectorAll('input').forEach(i => {
                i.required = false;
                i.value = '';
            });
        }
    }

    // Job Logic
    if (jobSection) {
        if (val === 'Job') {
            jobSection.classList.remove('hidden');
            jobSection.querySelectorAll('input, textarea').forEach(i => {
                if (i.name === 'company' || i.name === 'skills' || i.name === 'job_start_date' || i.name === 'pay_range' || i.name === 'job_apply_url') {
                    i.required = true;
                }
            });

            // Hide Image Upload for Jobs
            if (imageSection) imageSection.classList.add('hidden');
        } else {
            jobSection.classList.add('hidden');
            jobSection.querySelectorAll('input, textarea').forEach(i => {
                i.required = false;
                i.value = '';
            });

            // Show Image Upload for others
            if (imageSection) imageSection.classList.remove('hidden');
        }
    }

    // Address & Label Logic
    if (addressInput) {
        if (val === 'Job') {
            addressInput.required = true;
            addressInput.placeholder = "City, Country or Address";
            if (locationLabel) locationLabel.textContent = "Location";
            if (descriptionLabel) descriptionLabel.textContent = "Job Description";
        } else {
            if (descriptionLabel) descriptionLabel.textContent = "Description";
            // Address is optional for Service, Request, Product, and Event
            if (val === 'Service' || val === 'Request' || val === 'Product' || val === 'Event') {
                addressInput.required = false;
                addressInput.placeholder = "Address (Optional)";
                if (locationLabel) locationLabel.textContent = "Address (Optional)";
            } else {
                // Business, Food
                addressInput.required = true;
                addressInput.placeholder = "Start typing address...";
                if (locationLabel) locationLabel.textContent = "Address (Validated)";
            }
        }
    }

    // Hours of Operation Logic
    if (hoursSection) {
        const hoursInput = hoursSection.querySelector('input');
        // Allowed: Business, Service, Food
        if (val === 'Business' || val === 'Service' || val === 'Food') {
            hoursSection.classList.remove('hidden');
            hoursSection.style.display = ''; // Clear any inline styles
            if (hoursInput) hoursInput.disabled = false;
        } else {
            // Product, Event, Job, Request
            hoursSection.classList.add('hidden');
            hoursSection.style.display = ''; // Clear any inline styles
            if (hoursInput) {
                hoursInput.value = '';
                hoursInput.disabled = true;
            }
        }
    }
}
