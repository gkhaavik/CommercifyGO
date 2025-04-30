package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/postgres"
)

// RepositoryProvider provides access to all repositories
type RepositoryProvider interface {
	UserRepository() repository.UserRepository
	ProductRepository() repository.ProductRepository
	ProductVariantRepository() repository.ProductVariantRepository
	CategoryRepository() repository.CategoryRepository
	OrderRepository() repository.OrderRepository
	CartRepository() repository.CartRepository
	DiscountRepository() repository.DiscountRepository
	WebhookRepository() repository.WebhookRepository
}

// repositoryProvider is the concrete implementation of RepositoryProvider
type repositoryProvider struct {
	container Container
	mu        sync.Mutex

	userRepo           repository.UserRepository
	productRepo        repository.ProductRepository
	productVariantRepo repository.ProductVariantRepository
	categoryRepo       repository.CategoryRepository
	orderRepo          repository.OrderRepository
	cartRepo           repository.CartRepository
	discountRepo       repository.DiscountRepository
	webhookRepo        repository.WebhookRepository
}

// NewRepositoryProvider creates a new repository provider
func NewRepositoryProvider(container Container) RepositoryProvider {
	return &repositoryProvider{
		container: container,
	}
}

// UserRepository returns the user repository
func (p *repositoryProvider) UserRepository() repository.UserRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.userRepo == nil {
		p.userRepo = postgres.NewUserRepository(p.container.DB())
	}
	return p.userRepo
}

// ProductRepository returns the product repository
func (p *repositoryProvider) ProductRepository() repository.ProductRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productRepo == nil {
		p.productRepo = postgres.NewProductRepository(p.container.DB())
	}
	return p.productRepo
}

// ProductVariantRepository returns the product variant repository
func (p *repositoryProvider) ProductVariantRepository() repository.ProductVariantRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.productVariantRepo == nil {
		p.productVariantRepo = postgres.NewProductVariantRepository(p.container.DB())
	}
	return p.productVariantRepo
}

// CategoryRepository returns the category repository
func (p *repositoryProvider) CategoryRepository() repository.CategoryRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.categoryRepo == nil {
		p.categoryRepo = postgres.NewCategoryRepository(p.container.DB())
	}
	return p.categoryRepo
}

// OrderRepository returns the order repository
func (p *repositoryProvider) OrderRepository() repository.OrderRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.orderRepo == nil {
		p.orderRepo = postgres.NewOrderRepository(p.container.DB())
	}
	return p.orderRepo
}

// CartRepository returns the cart repository
func (p *repositoryProvider) CartRepository() repository.CartRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cartRepo == nil {
		p.cartRepo = postgres.NewCartRepository(p.container.DB())
	}
	return p.cartRepo
}

// DiscountRepository returns the discount repository
func (p *repositoryProvider) DiscountRepository() repository.DiscountRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.discountRepo == nil {
		p.discountRepo = postgres.NewDiscountRepository(p.container.DB())
	}
	return p.discountRepo
}

// WebhookRepository returns the webhook repository
func (p *repositoryProvider) WebhookRepository() repository.WebhookRepository {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.webhookRepo == nil {
		p.webhookRepo = postgres.NewWebhookRepository(p.container.DB())
	}
	return p.webhookRepo
}
