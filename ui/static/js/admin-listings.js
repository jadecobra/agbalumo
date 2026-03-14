function initAdminListings() {
    const selectAll = document.getElementById('selectAll');
    const checkboxes = document.querySelectorAll('.row-checkbox');
    const selectionCount = document.getElementById('selectionCount');
    const bulkActionSelect = document.getElementById('bulkActionSelect');
    const bulkActionButton = document.getElementById('bulkActionButton');

    const newCategorySelect = document.getElementById('newCategorySelect');

    function updateSelectionState() {
        const checkedCount = document.querySelectorAll('.row-checkbox:checked').length;
        if (selectionCount) {
            selectionCount.textContent = `${checkedCount} selected`;
        }

        const hasSelection = checkedCount > 0;
        const action = bulkActionSelect ? bulkActionSelect.value : "";
        
        if (bulkActionSelect) bulkActionSelect.disabled = !hasSelection;
        
        if (newCategorySelect) {
            if (action === "change_category" && hasSelection) {
                newCategorySelect.classList.remove('hidden');
            } else {
                newCategorySelect.classList.add('hidden');
            }
        }

        if (bulkActionButton) {
            let isDisabled = !hasSelection || action === "";
            if (action === "change_category" && newCategorySelect && newCategorySelect.value === "") {
                isDisabled = true;
            }
            bulkActionButton.disabled = isDisabled;
        }

        if (selectAll) {
            selectAll.checked = checkedCount === checkboxes.length && checkboxes.length > 0;
            selectAll.indeterminate = checkedCount > 0 && checkedCount < checkboxes.length;
        }
    }

    if (selectAll) {
        selectAll.addEventListener('change', function () {
            checkboxes.forEach(cb => {
                cb.checked = selectAll.checked;
            });
            updateSelectionState();
        });
    }

    checkboxes.forEach(cb => {
        cb.addEventListener('change', updateSelectionState);
    });

    if (bulkActionSelect) {
        bulkActionSelect.addEventListener('change', updateSelectionState);
    }

    if (newCategorySelect) {
        newCategorySelect.addEventListener('change', updateSelectionState);
    }

    updateSelectionState();
}

document.addEventListener('DOMContentLoaded', initAdminListings);
