package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/webhooks"
	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/container"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/handler"
	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
)

// Server represents the API server
type Server struct {
	config     *config.Config
	router     *mux.Router
	httpServer *http.Server
	logger     logger.Logger
	container  container.Container
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, db *sql.DB, logger logger.Logger) *Server {
	// Initialize dependency container
	diContainer := container.NewContainer(cfg, db, logger)

	// Post-initialization to break circular dependencies
	if cfg.MobilePay.Enabled {
		// Connect MobilePay service to WebhookService
		mobilePayService := diContainer.Services().MobilePayService()
		webhookService := diContainer.Services().WebhookService()
		if mobilePayService != nil && webhookService != nil {
			webhookService.SetMobilePayService(mobilePayService)
		}
	}

	router := mux.NewRouter()

	server := &Server{
		config:    cfg,
		router:    router,
		logger:    logger,
		container: diContainer,
	}

	// Apply CORS middleware to all routes
	// corsMiddleware := diContainer.Middlewares().CorsMiddleware()
	// router.Use(corsMiddleware.ApplyCors)

	server.setupRoutes()

	// Create HTTP server
	server.httpServer = &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return server
}

// setupRoutes configures all routes for the API
func (s *Server) setupRoutes() {
	// Extract handlers from container
	userHandler := s.container.Handlers().UserHandler()
	productHandler := s.container.Handlers().ProductHandler()
	cartHandler := s.container.Handlers().CartHandler()
	orderHandler := s.container.Handlers().OrderHandler()
	paymentHandler := s.container.Handlers().PaymentHandler()
	webhookHandler := s.container.Handlers().WebhookHandler()
	discountHandler := s.container.Handlers().DiscountHandler()
	shippingHandler := s.container.Handlers().ShippingHandler()
	currencyHandler := s.container.Handlers().CurrencyHandler()

	// Extract middleware from container
	authMiddleware := s.container.Middlewares().AuthMiddleware()

	// Register routes
	api := s.router.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/auth/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/signin", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/products/{productId:[0-9]+}", productHandler.GetProduct).Methods(http.MethodGet)

	api.HandleFunc("/products/search", productHandler.SearchProducts).Methods(http.MethodGet)
	api.HandleFunc("/categories", productHandler.ListCategories).Methods(http.MethodGet)
	api.HandleFunc("/payment/providers", paymentHandler.GetAvailablePaymentProviders).Methods(http.MethodGet)

	// Public discount routes
	api.HandleFunc("/discounts/validate", discountHandler.ValidateDiscountCode).Methods(http.MethodPost)
	api.HandleFunc("/guest/discounts/apply/{orderId:[0-9]+}", discountHandler.ApplyDiscountToGuestOrder).Methods(http.MethodPost)

	// Public currency routes
	api.HandleFunc("/currencies", currencyHandler.ListEnabledCurrencies).Methods(http.MethodGet)
	api.HandleFunc("/currencies/default", currencyHandler.GetDefaultCurrency).Methods(http.MethodGet)
	api.HandleFunc("/currencies/convert", currencyHandler.ConvertAmount).Methods(http.MethodPost)

	// Public shipping routes
	api.HandleFunc("/shipping/methods", shippingHandler.ListShippingMethods).Methods(http.MethodGet)
	api.HandleFunc("/shipping/methods/{shippingMethodId:[0-9]+}", shippingHandler.GetShippingMethodByID).Methods(http.MethodGet)
	api.HandleFunc("/shipping/options", shippingHandler.CalculateShippingOptions).Methods(http.MethodPost)
	api.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}/cost", shippingHandler.GetShippingCost).Methods(http.MethodPost)

	// Guest cart routes (no authentication required)
	api.HandleFunc("/guest/cart", cartHandler.GetCart).Methods(http.MethodGet)
	api.HandleFunc("/guest/cart/items", cartHandler.AddToCart).Methods(http.MethodPost)
	api.HandleFunc("/guest/cart/items/{productId:[0-9]+}", cartHandler.UpdateCartItem).Methods(http.MethodPut)
	api.HandleFunc("/guest/cart/items/{productId:[0-9]+}", cartHandler.RemoveFromCart).Methods(http.MethodDelete)
	api.HandleFunc("/guest/cart", cartHandler.ClearCart).Methods(http.MethodDelete)

	// Guest checkout route
	api.HandleFunc("/guest/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	api.HandleFunc("/guest/orders/{orderId:[0-9]+}/payment", orderHandler.ProcessPayment).Methods(http.MethodPost)

	// Convert guest cart to user cart after login
	api.HandleFunc("/guest/cart/convert", cartHandler.ConvertGuestCartToUserCart).Methods(http.MethodPost)

	// Webhooks
	api.HandleFunc("/webhooks/stripe", webhookHandler.HandleStripeWebhook).Methods(http.MethodPost)

	// Setup payment provider webhooks
	s.setupMobilePayWebhooks(api, webhookHandler)
	s.setupStripeWebhooks(api, webhookHandler)

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Authenticate)

	// User routes
	protected.HandleFunc("/users/me", userHandler.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/users/me/password", userHandler.ChangePassword).Methods(http.MethodPut)

	// Cart routes
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods(http.MethodGet)
	protected.HandleFunc("/cart/items", cartHandler.AddToCart).Methods(http.MethodPost)
	protected.HandleFunc("/cart/items/{productId:[0-9]+}", cartHandler.UpdateCartItem).Methods(http.MethodPut)
	protected.HandleFunc("/cart/items/{productId:[0-9]+}", cartHandler.RemoveFromCart).Methods(http.MethodDelete)
	protected.HandleFunc("/cart", cartHandler.ClearCart).Methods(http.MethodDelete)

	// Order routes
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	protected.HandleFunc("/orders/{orderId:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods(http.MethodGet)
	protected.HandleFunc("/orders/{orderId:[0-9]+}/payment", orderHandler.ProcessPayment).Methods(http.MethodPost)

	// Discount routes
	protected.HandleFunc("/discounts", discountHandler.CreateDiscount).Methods(http.MethodPost)
	protected.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.UpdateDiscount).Methods(http.MethodPut)
	protected.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.DeleteDiscount).Methods(http.MethodDelete)
	protected.HandleFunc("/discounts", discountHandler.ListDiscounts).Methods(http.MethodGet)
	protected.HandleFunc("/discounts/active", discountHandler.ListActiveDiscounts).Methods(http.MethodGet)
	protected.HandleFunc("/discounts/apply/{orderId:[0-9]+}", discountHandler.ApplyDiscountToOrder).Methods(http.MethodPost)
	protected.HandleFunc("/discounts/remove/{orderId:[0-9]+}", discountHandler.RemoveDiscountFromOrder).Methods(http.MethodDelete)
	protected.HandleFunc("/discounts/{discountId:[0-9]+}", discountHandler.GetDiscount).Methods(http.MethodGet)

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminOnly)
	admin.HandleFunc("/users", userHandler.ListUsers).Methods(http.MethodGet)
	admin.HandleFunc("/orders", orderHandler.ListAllOrders).Methods(http.MethodGet)
	admin.HandleFunc("/orders/{orderId:[0-9]+}/status", orderHandler.UpdateOrderStatus).Methods(http.MethodPut)

	// Admin currency routes
	admin.HandleFunc("/currencies/all", currencyHandler.ListCurrencies).Methods(http.MethodGet)
	admin.HandleFunc("/currencies", currencyHandler.CreateCurrency).Methods(http.MethodPost)
	admin.HandleFunc("/currencies", currencyHandler.UpdateCurrency).Methods(http.MethodPut)
	admin.HandleFunc("/currencies", currencyHandler.DeleteCurrency).Methods(http.MethodDelete)
	admin.HandleFunc("/currencies/default", currencyHandler.SetDefaultCurrency).Methods(http.MethodPut)

	// Shipping management routes (admin only)
	admin.HandleFunc("/shipping/methods", shippingHandler.CreateShippingMethod).Methods(http.MethodPost)
	admin.HandleFunc("/shipping/methods/{shippingMethodId:[0-9]+}", shippingHandler.UpdateShippingMethod).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/zones", shippingHandler.CreateShippingZone).Methods(http.MethodPost)
	admin.HandleFunc("/shipping/zones", shippingHandler.ListShippingZones).Methods(http.MethodGet)
	admin.HandleFunc("/shipping/zones/{shippingZoneId:[0-9]+}", shippingHandler.GetShippingZoneByID).Methods(http.MethodGet)
	admin.HandleFunc("/shipping/zones/{shippingZoneId:[0-9]+}", shippingHandler.UpdateShippingZone).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/rates", shippingHandler.CreateShippingRate).Methods(http.MethodPost)
	admin.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}", shippingHandler.GetShippingRateByID).Methods(http.MethodGet)
	admin.HandleFunc("/shipping/rates/{shippingRateId:[0-9]+}", shippingHandler.UpdateShippingRate).Methods(http.MethodPut)
	admin.HandleFunc("/shipping/rates/weight", shippingHandler.CreateWeightBasedRate).Methods(http.MethodPost)
	admin.HandleFunc("/shipping/rates/value", shippingHandler.CreateValueBasedRate).Methods(http.MethodPost)

	// Payment management routes (admin only)
	admin.HandleFunc("/payments/{paymentId}/capture", paymentHandler.CapturePayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/cancel", paymentHandler.CancelPayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/refund", paymentHandler.RefundPayment).Methods(http.MethodPost)
	admin.HandleFunc("/payments/{paymentId}/force-approve", paymentHandler.ForceApproveMobilePayPayment).Methods(http.MethodPost)

	// Webhook management routes (admin only)
	admin.HandleFunc("/webhooks", webhookHandler.ListWebhooks).Methods(http.MethodGet)
	admin.HandleFunc("/webhooks/{webhookId:[0-9]+}", webhookHandler.GetWebhook).Methods(http.MethodGet)
	admin.HandleFunc("/webhooks/{webhookId:[0-9]+}", webhookHandler.DeleteWebhook).Methods(http.MethodDelete)
	admin.HandleFunc("/webhooks/mobilepay", webhookHandler.RegisterMobilePayWebhook).Methods(http.MethodPost)
	admin.HandleFunc("/webhooks/mobilepay", webhookHandler.GetMobilePayWebhooks).Methods(http.MethodGet)

	admin.HandleFunc("/products", productHandler.ListProducts).Methods(http.MethodGet)
	admin.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	admin.HandleFunc("/products/{productId:[0-9]+}", productHandler.UpdateProduct).Methods(http.MethodPut)
	admin.HandleFunc("/products/{productId:[0-9]+}", productHandler.DeleteProduct).Methods(http.MethodDelete)

	// Product variant routes
	admin.HandleFunc("/products/{productId:[0-9]+}/variants", productHandler.AddVariant).Methods(http.MethodPost)
	admin.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.UpdateVariant).Methods(http.MethodPut)
	admin.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.DeleteVariant).Methods(http.MethodDelete)
}

// setupStripeWebhooks configures Stripe webhooks
func (s *Server) setupStripeWebhooks(api *mux.Router, webhookHandler *handler.WebhookHandler) {
	if !s.config.Stripe.Enabled {
		return
	}

	if s.config.Stripe.WebhookSecret == "" {
		s.logger.Warn("Stripe webhook secret is not configured, webhooks will not validate signatures")
	} else {
		s.logger.Info("Stripe webhook endpoint configured at /api/webhooks/stripe")
	}

	// Note: For Stripe, webhook endpoints are already registered in the routes.
	// We don't need to dynamically register them like in MobilePay.
	// This method exists for consistency with MobilePay setup and to handle any future
	// Stripe webhook configuration needs.
}

// setupMobilePayWebhooks configures MobilePay webhooks if enabled
func (s *Server) setupMobilePayWebhooks(api *mux.Router, webhookHandler *handler.WebhookHandler) {
	if !s.config.MobilePay.Enabled {
		return
	}

	// Get webhooks
	webhookUseCase := s.container.UseCases().WebhookUseCase()
	result, err := webhookUseCase.GetAllWebhooks()
	if err != nil {
		s.logger.Error("Failed to get MobilePay webhooks: %v", err)
		return
	}

	// Register webhook if none exists
	if len(result) == 0 {
		webhookService := s.container.Services().WebhookService()
		webhook, err := webhookService.RegisterMobilePayWebhook(s.config.MobilePay.WebhookURL, []string{
			string(models.WebhookEventPaymentAborted),
			string(models.WebhookEventPaymentCancelled),
			string(models.WebhookEventPaymentCaptured),
			string(models.WebhookEventPaymentRefunded),
			string(models.WebhookEventPaymentExpired),
			string(models.WebhookEventPaymentAuthorized),
		})

		if err != nil {
			s.logger.Error("Failed to register MobilePay webhook: %v", err)
		} else {
			s.logger.Info("Registered new MobilePay webhook: %s", webhook.URL)
			result = append(result, webhook)
		}
	} else {
		s.logger.Info("Found %d MobilePay webhooks", len(result))
	}

	// Configure webhook handlers
	for _, webhook := range result {
		if webhook.IsActive && webhook.Provider == "mobilepay" {
			handler := webhooks.NewHandler(webhook.Secret)
			router := webhooks.NewRouter()

			router.HandleFunc(models.EventAuthorized, webhookHandler.HandleMobilePayAuthorized)
			router.HandleFunc(models.EventAborted, webhookHandler.HandleMobilePayAborted)
			router.HandleFunc(models.EventCancelled, webhookHandler.HandleMobilePayCancelled)
			router.HandleFunc(models.EventCaptured, webhookHandler.HandleMobilePayCaptured)
			router.HandleFunc(models.EventRefunded, webhookHandler.HandleMobilePayRefunded)
			router.HandleFunc(models.EventExpired, webhookHandler.HandleMobilePayExpired)

			router.HandleDefault(func(event *models.WebhookEvent) error {
				fmt.Printf("Received unhandled event: %s\n", event.Name)
				return nil
			})

			api.HandleFunc("/webhooks/mobilepay", handler.HandleHTTP(router.Process))
			s.logger.Info("Registered MobilePay webhook: %s", webhook.URL)
		}
	}
}

// Start starts the server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
