package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/application/usecase"
)

// UseCaseProvider provides access to all use cases
type UseCaseProvider interface {
	UserUseCase() *usecase.UserUseCase
	ProductUseCase() *usecase.ProductUseCase
	CartUseCase() *usecase.CartUseCase
	OrderUseCase() *usecase.OrderUseCase
	DiscountUseCase() *usecase.DiscountUseCase
	WebhookUseCase() *usecase.WebhookUseCase
}

// useCaseProvider is the concrete implementation of UseCaseProvider
type useCaseProvider struct {
	container Container
	mu        sync.Mutex

	userUseCase     *usecase.UserUseCase
	productUseCase  *usecase.ProductUseCase
	cartUseCase     *usecase.CartUseCase
	orderUseCase    *usecase.OrderUseCase
	discountUseCase *usecase.DiscountUseCase
	webhookUseCase  *usecase.WebhookUseCase
}

// NewUseCaseProvider creates a new use case provider
func NewUseCaseProvider(container Container) UseCaseProvider {
	return &useCaseProvider{
		container: container,
	}
}

// UserUseCase returns the user use case
func (p *useCaseProvider) UserUseCase() *usecase.UserUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userUseCase == nil {
		p.userUseCase = usecase.NewUserUseCase(
			p.container.Repositories().UserRepository(),
		)
	}
	return p.userUseCase
}

// ProductUseCase returns the product use case
func (p *useCaseProvider) ProductUseCase() *usecase.ProductUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productUseCase == nil {
		p.productUseCase = usecase.NewProductUseCase(
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().CategoryRepository(),
			p.container.Repositories().ProductVariantRepository(),
		)
	}
	return p.productUseCase
}

// CartUseCase returns the cart use case
func (p *useCaseProvider) CartUseCase() *usecase.CartUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cartUseCase == nil {
		p.cartUseCase = usecase.NewCartUseCase(
			p.container.Repositories().CartRepository(),
			p.container.Repositories().ProductRepository(),
		)
	}
	return p.cartUseCase
}

// OrderUseCase returns the order use case
func (p *useCaseProvider) OrderUseCase() *usecase.OrderUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.orderUseCase == nil {
		p.orderUseCase = usecase.NewOrderUseCase(
			p.container.Repositories().OrderRepository(),
			p.container.Repositories().CartRepository(),
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().UserRepository(),
			p.container.Services().PaymentService(),
			p.container.Services().EmailService(),
		)
	}
	return p.orderUseCase
}

// DiscountUseCase returns the discount use case
func (p *useCaseProvider) DiscountUseCase() *usecase.DiscountUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discountUseCase == nil {
		p.discountUseCase = usecase.NewDiscountUseCase(
			p.container.Repositories().DiscountRepository(),
			p.container.Repositories().ProductRepository(),
			p.container.Repositories().CategoryRepository(),
			p.container.Repositories().OrderRepository(),
		)
	}
	return p.discountUseCase
}

// WebhookUseCase returns the webhook use case
func (p *useCaseProvider) WebhookUseCase() *usecase.WebhookUseCase {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.webhookUseCase == nil {
		p.webhookUseCase = usecase.NewWebhookUseCase(
			p.container.Repositories().WebhookRepository(),
			p.container.Services().WebhookService(),
		)
	}
	return p.webhookUseCase
}
