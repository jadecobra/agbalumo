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
                
                var removeInput = document.getElementById('remove-image-' + listingId);
                if (removeInput) {
                    removeInput.value = 'false';
                }
            }
        });
    }
}

// Clear existing image
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

// Remove new image preview
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
    if (existing && (!removeInput || removeInput.value !== 'true')) {
        existing.classList.remove('hidden');
    }
}
