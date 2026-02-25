function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.showModal();
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.close();
        modal.remove();
    }
}

document.addEventListener('DOMContentLoaded', function () {
    document.querySelectorAll('dialog[id][data-auto-open]').forEach(modal => {
        modal.showModal();
    });

    document.querySelectorAll('dialog[id]').forEach(modal => {
        modal.addEventListener('click', function (event) {
            const rect = modal.getBoundingClientRect();
            const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height
                && rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
            if (!isInDialog) {
                modal.close();
            }
        });
    });

    // Auto-open detail modals
    document.querySelectorAll('dialog[id^="detail-modal-"]').forEach(modal => {
        modal.showModal();
    });

    // Auto-open feedback modal
    const feedbackModal = document.getElementById('feedback-modal');
    if (feedbackModal) {
        feedbackModal.showModal();
    }

    // Initialize edit modals if present on initial load
    document.querySelectorAll('dialog[id^="edit-listing-modal-"]').forEach(modal => {
        const listingId = modal.id.replace('edit-listing-modal-', '');
        if (listingId) {
            initEditTypeToggle(listingId);
            initEditMaps(listingId);
            initEditImagePreview(listingId);
        }
    });
});

// HTMX Listener for dynamically loaded content
document.addEventListener('htmx:afterSettle', function (event) {
    const target = event.detail.target;
    if (!target) return;

    // Handle new edit listing modals
    const editModal = target.querySelector('dialog[id^="edit-listing-modal-"]') ||
        (target.id?.startsWith('edit-listing-modal-') ? target : null);

    if (editModal) {
        const listingId = editModal.id.replace('edit-listing-modal-', '');
        if (listingId) {
            editModal.showModal();
            initEditTypeToggle(listingId);
            initEditMaps(listingId);
            initEditImagePreview(listingId);
        }
    }

    // Auto-open detail modals loaded via HTMX
    target.querySelectorAll('dialog[id^="detail-modal-"]').forEach(modal => {
        modal.showModal();
    });

    // Auto-open feedback modal loaded via HTMX
    const feedbackModal = target.querySelector('#feedback-modal') ||
        (target.id === 'feedback-modal' ? target : null);
    if (feedbackModal) {
        feedbackModal.showModal();
    }
});

// Image preview for create listing
function initCreateImagePreview() {
    var imageInput = document.getElementById('image-input');
    var previewContainer = document.getElementById('image-preview-container');
    var previewImg = document.getElementById('image-preview');
    var removeBtn = document.getElementById('remove-image-btn');

    if (imageInput) {
        imageInput.addEventListener('change', function (e) {
            var file = e.target.files[0];
            if (file) {
                var reader = new FileReader();
                reader.onload = function (e) {
                    previewImg.src = e.target.result;
                    previewContainer.classList.remove('hidden');
                };
                reader.readAsDataURL(file);
            }
        });
    }

    if (removeBtn) {
        removeBtn.addEventListener('click', function () {
            imageInput.value = '';
            previewContainer.classList.add('hidden');
        });
    }
}

// Type field toggle for edit listing
function initEditTypeToggle(listingId) {
    const form = document.querySelector('#edit-listing-modal-' + listingId + ' form');
    if (!form) return;

    const typeSelect = form.querySelector('select[name="type"]');
    const eventSection = document.getElementById('edit-event-dates-section-' + listingId);
    const requestSection = document.getElementById('edit-request-section-' + listingId);
    const jobSection = document.getElementById('edit-job-section-' + listingId);

    function toggleFields() {
        const val = typeSelect.value;

        eventSection.classList.add('hidden');
        requestSection.classList.add('hidden');
        jobSection.classList.add('hidden');

        [eventSection, requestSection, jobSection].forEach(section => {
            section.querySelectorAll('input').forEach(i => i.required = false);
        });

        if (val === 'Event') {
            eventSection.classList.remove('hidden');
            eventSection.querySelectorAll('input').forEach(i => i.required = true);
        } else if (val === 'Request') {
            requestSection.classList.remove('hidden');
            requestSection.querySelectorAll('input').forEach(i => i.required = true);
        } else if (val === 'Job') {
            jobSection.classList.remove('hidden');
            jobSection.querySelectorAll('input').forEach(i => i.required = true);
        }

        // Address Requirement Logic
        const addressInput = document.getElementById('edit-address-input-' + listingId);
        if (addressInput) {
            if (val === 'Service' || val === 'Job' || val === 'Request') {
                addressInput.required = false;
                addressInput.placeholder = "Address (Optional)";
                const label = addressInput.previousElementSibling;
                if (label) label.textContent = "Address (Optional)";
            } else {
                addressInput.required = true;
                addressInput.placeholder = "Start typing address...";
                const label = addressInput.previousElementSibling;
                if (label) label.textContent = "Address (Validated)";
            }
        }

        // Hours of Operation Logic
        const hoursInput = document.querySelector('#edit-listing-modal-' + listingId + ' input[name="hours_of_operation"]');
        if (hoursInput) {
            const hoursContainer = hoursInput.closest('div');
            // Allowed: Business, Service, Food
            if (val === 'Business' || val === 'Service' || val === 'Food') {
                if (hoursContainer) hoursContainer.classList.remove('hidden');
                hoursInput.disabled = false;
            } else {
                if (hoursContainer) hoursContainer.classList.add('hidden');
                hoursInput.value = '';
                hoursInput.disabled = true;
            }
        }
    }

    if (typeSelect) {
        typeSelect.addEventListener('change', toggleFields);
        toggleFields();
    }
}

// Google Maps autocomplete for edit listing
function initEditMaps(listingId) {
    const input = document.getElementById('edit-address-input-' + listingId);
    if (!input) return;

    if (typeof google === 'undefined' || !google.maps || !google.maps.places) return;

    const options = {
        fields: ["address_components", "geometry", "formatted_address"],
        types: ["address"],
    };
    const autocomplete = new google.maps.places.Autocomplete(input, options);
    autocomplete.addListener("place_changed", () => {
        const place = autocomplete.getPlace();
        let city = "";
        for (const component of place.address_components) {
            const types = component.types;
            if (types.includes("locality")) {
                city = component.long_name;
                break;
            }
            if (types.includes("postal_town")) {
                city = component.long_name;
            }
            if (!city && types.includes("administrative_area_level_2")) {
                city = component.long_name;
            }
        }
        const cityInput = document.getElementById('edit-city-input-' + listingId);
        if (cityInput) {
            cityInput.value = city || "Unknown";
        }
    });
}

// Image preview for edit listing
function initEditImagePreview(listingId) {
    var imageInput = document.getElementById('edit-image-input-' + listingId);
    var previewContainer = document.getElementById('edit-image-preview-container-' + listingId);
    var previewImg = document.getElementById('edit-image-preview-' + listingId);
    var existingImage = document.getElementById('edit-existing-image-' + listingId);

    if (imageInput) {
        imageInput.addEventListener('change', function (e) {
            var file = e.target.files[0];
            if (file) {
                var reader = new FileReader();
                reader.onload = function (e) {
                    previewImg.src = e.target.result;
                    previewContainer.classList.remove('hidden');
                    if (existingImage) existingImage.classList.add('hidden');
                };
                reader.readAsDataURL(file);
            }
        });
    }
}
// Clear existing image (for database update)
function clearEditImage(listingId) {
    const removeInput = document.getElementById('remove-image-' + listingId);
    if (removeInput) {
        removeInput.value = 'true';
    }
    const container = document.getElementById('edit-existing-image-' + listingId);
    if (container) {
        container.classList.add('hidden');
    }
}

// Remove new image preview (just purely cosmetic)
function removeEditImagePreview(listingId) {
    const input = document.getElementById('edit-image-input-' + listingId);
    if (input) {
        input.value = '';
    }
    const container = document.getElementById('edit-image-preview-container-' + listingId);
    if (container) {
        container.classList.add('hidden');
    }
    const existing = document.getElementById('edit-existing-image-' + listingId);
    const removeInput = document.getElementById('remove-image-' + listingId);
    // If we had an existing image that was hidden by the preview, show it again
    // unless it was explicitly removed.
    if (existing && (!removeInput || removeInput.value !== 'true')) {
        existing.classList.remove('hidden');
    }
}
