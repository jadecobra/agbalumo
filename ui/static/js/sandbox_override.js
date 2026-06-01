// Prevent native dialog elements from entering the Top Layer in the sandbox environment
if (window.HTMLDialogElement) {
    window.HTMLDialogElement.prototype.showModal = function() {
        this.setAttribute('open', '');
    };
    window.HTMLDialogElement.prototype.close = function() {
        // Keep them open permanently for visual inspection in the sandbox
    };
}
