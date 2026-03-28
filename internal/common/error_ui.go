package common

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func IsImageError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "File size exceeds") ||
		strings.Contains(msg, "Invalid file type") ||
		strings.Contains(msg, "Invalid or unsupported image")
}

func RenderImageErrorToast(c echo.Context, err error) error {
	var msg string
	if he, ok := err.(*echo.HTTPError); ok {
		if m, ok := he.Message.(string); ok {
			msg = m
		}
	}

	toastID := uuid.New().String()

	c.Response().Header().Set("HX-Reswap", "none")
	c.Response().Header().Set("Content-Type", "text/html")

	// #nosec - toastID and msg are manually escaped below
	return c.HTML(http.StatusBadRequest, fmt.Sprintf(`
	<div id="toast-%s" 
	     class="fixed top-4 right-4 z-50 max-w-sm w-full bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-xl shadow-lg p-4 flex items-start gap-3 animate-in slide-in-from-top-2 fade-in"
	     role="alert"
	     hx-on::after-transaction="if(event.detail.failed) setTimeout(() => { const t = document.getElementById('toast-%s'); if(t) { t.style.animation = 'fade-out 0.3s ease-out forwards'; setTimeout(() => t.remove(), 300); } }, 5000)">
	    <span class="material-symbols-outlined text-red-500 text-[20px] mt-0.5">error</span>
	    <div class="flex-1 min-w-0">
	        <p class="text-sm font-medium text-red-800 dark:text-red-200">Image Upload Failed</p>
	        <p class="text-sm text-red-600 dark:text-red-300 mt-1">%s</p>
	    </div>
	    <button hx-on:click="this.parentElement.remove()" 
	            class="text-red-400 hover:text-red-600 dark:hover:text-red-200 transition-colors">
	        <span class="material-symbols-outlined text-[18px]">close</span>
	    </button>
	</div>`, 
	template.HTMLEscapeString(toastID), 
	template.HTMLEscapeString(toastID), 
	template.HTMLEscapeString(msg)))
}
