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

    // Auto-open any loaded data-auto-open modals
    const autoOpenModals = target.querySelectorAll('dialog[id][data-auto-open]');
    if (autoOpenModals) {
        autoOpenModals.forEach(modal => {
            modal.showModal();
        });
    }
    if (target.matches && target.matches('dialog[id][data-auto-open]')) {
        target.showModal();
    }
});

// ========== Image Upload Functions ==========

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
                showImageUploadLoading(imageInput, true);
                
                var reader = new FileReader();
                reader.onload = function (e) {
                    previewImg.src = e.target.result;
                    previewContainer.classList.remove('hidden');
                    showImageUploadLoading(imageInput, false);
                };
                reader.onerror = function() {
                    showImageUploadLoading(imageInput, false);
                    showImageError('Failed to read file');
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

// Show/hide loading state on image input (safe: no innerHTML)
function showImageUploadLoading(inputEl, show) {
    var form = inputEl.closest('form');
    if (!form) return;
    
    var loadingIndicator = form.querySelector('.image-upload-loading');
    if (show) {
        if (!loadingIndicator) {
            loadingIndicator = document.createElement('div');
            loadingIndicator.className = 'image-upload-loading flex items-center gap-2 text-sm text-stone-500 mt-1';
            var spinner = document.createElement('span');
            spinner.className = 'material-symbols-outlined animate-spin text-[16px]';
            spinner.textContent = 'progress_activity';
            var text = document.createTextNode(' Processing image...');
            loadingIndicator.appendChild(spinner);
            loadingIndicator.appendChild(text);
            inputEl.parentNode.insertBefore(loadingIndicator, inputEl.nextSibling);
        }
        loadingIndicator.classList.remove('hidden');
    } else if (loadingIndicator) {
        loadingIndicator.classList.add('hidden');
    }
}

// Show image upload error (safe: no innerHTML)
function showImageError(message) {
    var existingError = document.querySelector('.image-upload-error');
    if (existingError) {
        existingError.remove();
    }
    
    var errorDiv = document.createElement('div');
    errorDiv.className = 'image-upload-error flex items-center gap-2 text-sm text-red-600 dark:text-red-400 mt-1 px-3 py-2 bg-red-50 dark:bg-red-900/20 rounded-lg';
    var icon = document.createElement('span');
    icon.className = 'material-symbols-outlined text-[16px]';
    icon.textContent = 'error';
    var text = document.createTextNode(' ' + message);
    errorDiv.appendChild(icon);
    errorDiv.appendChild(text);
    
    var imageInput = document.querySelector('[name="image"]');
    if (imageInput) {
        var container = imageInput.closest('.flex.flex-col.gap-1\\.5');
        if (container) {
            container.appendChild(errorDiv);
        }
    }
    
    setTimeout(function() {
        errorDiv.remove();
    }, 5000);
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
                showImageUploadLoading(imageInput, true);
                
                var reader = new FileReader();
                reader.onload = function (e) {
                    previewImg.src = e.target.result;
                    previewContainer.classList.remove('hidden');
                    if (existingImage) existingImage.classList.add('hidden');
                    showImageUploadLoading(imageInput, false);
                };
                reader.onerror = function() {
                    showImageUploadLoading(imageInput, false);
                    showImageError('Failed to read file');
                };
                reader.readAsDataURL(file);
                
                // Reset remove_image flag since user selected a new image
                var removeInput = document.getElementById('remove-image-' + listingId);
                if (removeInput) {
                    removeInput.value = 'false';
                }
            }
        });
    }
}

// Clear existing image (for database update)
function clearEditImage(listingId) {
    var removeInput = document.getElementById('remove-image-' + listingId);
    if (removeInput) {
        removeInput.value = 'true';
    }
    var container = document.getElementById('edit-existing-image-' + listingId);
    if (container) {
        container.classList.add('hidden');
    }
}

// Remove new image preview (just purely cosmetic)
function removeEditImagePreview(listingId) {
    var input = document.getElementById('edit-image-input-' + listingId);
    if (input) {
        input.value = '';
    }
    var container = document.getElementById('edit-image-preview-container-' + listingId);
    if (container) {
        container.classList.add('hidden');
    }
    var existing = document.getElementById('edit-existing-image-' + listingId);
    var removeInput = document.getElementById('remove-image-' + listingId);
    // If we had an existing image that was hidden by the preview, show it again
    // unless it was explicitly removed.
    if (existing && (!removeInput || removeInput.value !== 'true')) {
        existing.classList.remove('hidden');
    }
}

// ========== Type & Fields Functions ==========

// Type field toggle for edit listing
function initEditTypeToggle(listingId) {
    var form = document.querySelector('#edit-listing-modal-' + listingId + ' form');
    if (!form) return;

    var typeSelect = form.querySelector('select[name="type"]');
    var eventSection = document.getElementById('edit-event-dates-section-' + listingId);
    var requestSection = document.getElementById('edit-request-section-' + listingId);
    var jobSection = document.getElementById('edit-job-section-' + listingId);

    function toggleFields() {
        var val = typeSelect.value;

        eventSection.classList.add('hidden');
        requestSection.classList.add('hidden');
        jobSection.classList.add('hidden');

        [eventSection, requestSection, jobSection].forEach(function(section) {
            section.querySelectorAll('input').forEach(function(i) { i.required = false; });
        });

        if (val === 'Event') {
            eventSection.classList.remove('hidden');
            eventSection.querySelectorAll('input').forEach(function(i) { i.required = true; });
        } else if (val === 'Request') {
            requestSection.classList.remove('hidden');
            requestSection.querySelectorAll('input').forEach(function(i) { i.required = true; });
        } else if (val === 'Job') {
            jobSection.classList.remove('hidden');
            jobSection.querySelectorAll('input').forEach(function(i) { i.required = true; });
        }

        // Address Requirement Logic
        var addressInput = document.getElementById('edit-address-input-' + listingId);
        if (addressInput) {
            if (val === 'Service' || val === 'Job' || val === 'Request') {
                addressInput.required = false;
                addressInput.placeholder = "Address (Optional)";
                var label = addressInput.previousElementSibling;
                if (label) label.textContent = "Address (Optional)";
            } else {
                addressInput.required = true;
                addressInput.placeholder = "Start typing address...";
                var label = addressInput.previousElementSibling;
                if (label) label.textContent = "Address (Validated)";
            }
        }

        // Hours of Operation Logic
        var hoursInput = document.querySelector('#edit-listing-modal-' + listingId + ' input[name="hours_of_operation"]');
        if (hoursInput) {
            var hoursContainer = hoursInput.closest('div');
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
    var input = document.getElementById('edit-address-input-' + listingId);
    if (!input) return;

    if (typeof google === 'undefined' || !google.maps || !google.maps.places) return;

    var options = {
        fields: ["address_components", "geometry", "formatted_address"],
        types: ["address"],
    };
    var autocomplete = new google.maps.places.Autocomplete(input, options);
    autocomplete.addEventListener("place_changed", function() {
        var place = autocomplete.getPlace();
        var city = "";
        for (var i = 0; i < place.address_components.length; i++) {
            var component = place.address_components[i];
            var types = component.types;
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
        var cityInput = document.getElementById('edit-city-input-' + listingId);
        if (cityInput) {
            cityInput.value = city || "Unknown";
        }
    });
}
