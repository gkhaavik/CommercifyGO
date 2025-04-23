package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/api/handler"
	"github.com/zenfulcode/commercify/internal/api/middleware"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/email"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/postgres"
)

// Server represents the API server
type Server struct {
	config     *config.Config
	router     *mux.Router
	httpServer *http.Server
	logger     logger.Logger
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, db *sql.DB, logger logger.Logger) *Server {
	router := mux.NewRouter()

	// Create repositories
	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	productVariantRepo := postgres.NewProductVariantRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	discountRepo := postgres.NewDiscountRepository(db)

	// Create services
	jwtService := auth.NewJWTService(cfg.Auth)

	// Create payment service with multiple providers
	paymentService := payment.NewMultiProviderPaymentService(cfg, logger)

	emailService := email.NewSMTPEmailService(cfg.Email, logger)

	// Create use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, productVariantRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, cartRepo, productRepo, userRepo, paymentService, emailService)
	discountUseCase := usecase.NewDiscountUseCase(discountRepo, productRepo, categoryRepo, orderRepo)

	// Create handlers
	userHandler := handler.NewUserHandler(userUseCase, jwtService, logger)
	productHandler := handler.NewProductHandler(productUseCase, logger)
	cartHandler := handler.NewCartHandler(cartUseCase, logger)
	orderHandler := handler.NewOrderHandler(orderUseCase, logger)
	paymentHandler := handler.NewPaymentHandler(orderUseCase, logger)
	webhookHandler := handler.NewWebhookHandler(cfg, orderUseCase, logger)
	discountHandler := handler.NewDiscountHandler(discountUseCase, orderUseCase, logger)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, logger)

	// Register routes
	api := router.PathPrefix("/api").Subrouter()

	// Public routes
	// api.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)
	api.HandleFunc("/users/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/users/login", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/products", productHandler.ListProducts).Methods(http.MethodGet)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProduct).Methods(http.MethodGet)
	api.HandleFunc("/products/search", productHandler.SearchProducts).Methods(http.MethodGet)
	api.HandleFunc("/categories", productHandler.ListCategories).Methods(http.MethodGet)
	api.HandleFunc("/payment/providers", paymentHandler.GetAvailablePaymentProviders).Methods(http.MethodGet)
	api.HandleFunc("/discounts/validate", discountHandler.ValidateDiscountCode).Methods(http.MethodPost)

	// Webhooks
	api.HandleFunc("/webhooks/stripe", webhookHandler.HandleStripeWebhook).Methods(http.MethodPost)

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.Authenticate)

	// User routes
	protected.HandleFunc("/users/me", userHandler.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods(http.MethodPut)
	protected.HandleFunc("/users/me/password", userHandler.ChangePassword).Methods(http.MethodPut)

	// Product routes (seller only)
	protected.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	protected.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct).Methods(http.MethodPut)
	protected.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct).Methods(http.MethodDelete)
	protected.HandleFunc("/products/seller", productHandler.ListSellerProducts).Methods(http.MethodGet)

	// Product variant routes (seller only)
	protected.HandleFunc("/products/{productId:[0-9]+}/variants", productHandler.AddVariant).Methods(http.MethodPost)
	protected.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.UpdateVariant).Methods(http.MethodPut)
	protected.HandleFunc("/products/{productId:[0-9]+}/variants/{variantId:[0-9]+}", productHandler.DeleteVariant).Methods(http.MethodDelete)

	// Cart routes
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods(http.MethodGet)
	protected.HandleFunc("/cart/items", cartHandler.AddToCart).Methods(http.MethodPost)
	protected.HandleFunc("/cart/items/{productId:[0-9]+}", cartHandler.UpdateCartItem).Methods(http.MethodPut)
	protected.HandleFunc("/cart/items/{productId:[0-9]+}", cartHandler.RemoveFromCart).Methods(http.MethodDelete)
	protected.HandleFunc("/cart", cartHandler.ClearCart).Methods(http.MethodDelete)

	// Order routes
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	protected.HandleFunc("/orders/{id:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods(http.MethodGet)
	protected.HandleFunc("/orders/{id:[0-9]+}/payment", orderHandler.ProcessPayment).Methods(http.MethodPost)

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
	admin.HandleFunc("/orders/{id:[0-9]+}/status", orderHandler.UpdateOrderStatus).Methods(http.MethodPut)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		config:     cfg,
		router:     router,
		httpServer: httpServer,
		logger:     logger,
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
