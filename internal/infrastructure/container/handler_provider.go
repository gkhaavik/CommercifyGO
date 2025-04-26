package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/interfaces/api/handler"
)

// HandlerProvider provides access to all handlers
type HandlerProvider interface {
	UserHandler() *handler.UserHandler
	ProductHandler() *handler.ProductHandler
	CartHandler() *handler.CartHandler
	OrderHandler() *handler.OrderHandler
	PaymentHandler() *handler.PaymentHandler
	WebhookHandler() *handler.WebhookHandler
	DiscountHandler() *handler.DiscountHandler
}

// handlerProvider is the concrete implementation of HandlerProvider
type handlerProvider struct {
	container Container
	mu        sync.Mutex

	userHandler     *handler.UserHandler
	productHandler  *handler.ProductHandler
	cartHandler     *handler.CartHandler
	orderHandler    *handler.OrderHandler
	paymentHandler  *handler.PaymentHandler
	webhookHandler  *handler.WebhookHandler
	discountHandler *handler.DiscountHandler
}

// NewHandlerProvider creates a new handler provider
func NewHandlerProvider(container Container) HandlerProvider {
	return &handlerProvider{
		container: container,
	}
}

// UserHandler returns the user handler
func (p *handlerProvider) UserHandler() *handler.UserHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userHandler == nil {
		p.userHandler = handler.NewUserHandler(
			p.container.UseCases().UserUseCase(),
			p.container.Services().JWTService(),
			p.container.Logger(),
		)
	}
	return p.userHandler
}

// ProductHandler returns the product handler
func (p *handlerProvider) ProductHandler() *handler.ProductHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productHandler == nil {
		p.productHandler = handler.NewProductHandler(
			p.container.UseCases().ProductUseCase(),
			p.container.Logger(),
		)
	}
	return p.productHandler
}

// CartHandler returns the cart handler
func (p *handlerProvider) CartHandler() *handler.CartHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cartHandler == nil {
		p.cartHandler = handler.NewCartHandler(
			p.container.UseCases().CartUseCase(),
			p.container.Logger(),
		)
	}
	return p.cartHandler
}

// OrderHandler returns the order handler
func (p *handlerProvider) OrderHandler() *handler.OrderHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.orderHandler == nil {
		p.orderHandler = handler.NewOrderHandler(
			p.container.UseCases().OrderUseCase(),
			p.container.Logger(),
		)
	}
	return p.orderHandler
}

// PaymentHandler returns the payment handler
func (p *handlerProvider) PaymentHandler() *handler.PaymentHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.paymentHandler == nil {
		p.paymentHandler = handler.NewPaymentHandler(
			p.container.UseCases().OrderUseCase(),
			p.container.Logger(),
		)
	}
	return p.paymentHandler
}

// WebhookHandler returns the webhook handler
func (p *handlerProvider) WebhookHandler() *handler.WebhookHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.webhookHandler == nil {
		p.webhookHandler = handler.NewWebhookHandler(
			p.container.Config(),
			p.container.UseCases().OrderUseCase(),
			p.container.UseCases().WebhookUseCase(),
			p.container.Logger(),
		)
	}
	return p.webhookHandler
}

// DiscountHandler returns the discount handler
func (p *handlerProvider) DiscountHandler() *handler.DiscountHandler {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discountHandler == nil {
		p.discountHandler = handler.NewDiscountHandler(
			p.container.UseCases().DiscountUseCase(),
			p.container.UseCases().OrderUseCase(),
			p.container.Logger(),
		)
	}
	return p.discountHandler
}
