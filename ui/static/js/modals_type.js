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
