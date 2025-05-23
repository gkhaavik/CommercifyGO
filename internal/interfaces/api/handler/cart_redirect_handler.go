package handler

import (
	"net/http"
	"strings"

	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// CartRedirectHandler handles requests to deprecated cart endpoints
// and redirects them to the equivalent checkout endpoints
type CartRedirectHandler struct {
	logger logger.Logger
}

// NewCartRedirectHandler creates a new cart redirect handler
func NewCartRedirectHandler(logger logger.Logger) *CartRedirectHandler {
	return &CartRedirectHandler{
		logger: logger,
	}
}

// RedirectToCheckout redirects cart requests to checkout endpoints
func (h *CartRedirectHandler) RedirectToCheckout(w http.ResponseWriter, r *http.Request) {
	// Log the deprecated endpoint use
	h.logger.Warn("Deprecated cart endpoint accessed: %s", r.URL.Path)

	// Map the cart URL to the equivalent checkout URL
	checkoutURL := strings.ReplaceAll(r.URL.Path, "/cart", "/checkout")

	// Create the redirect URL with the same query parameters
	redirectURL := checkoutURL
	if r.URL.RawQuery != "" {
		redirectURL += "?" + r.URL.RawQuery
	}

	// Send a 301 permanent redirect to the equivalent checkout endpoint
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}
