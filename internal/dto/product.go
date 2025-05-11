package dto

// ProductDTO represents a product in the system
type ProductDTO struct {
	BaseDTO
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	SKU           string       `json:"sku"`
	Price         float64      `json:"price"`
	StockQuantity int          `json:"stock_quantity"`
	Weight        float64      `json:"weight"`
	CategoryID    uint         `json:"category_id"`
	SellerID      uint         `json:"seller_id"`
	Images        []string     `json:"images"`
	HasVariants   bool         `json:"has_variants"`
	Variants      []VariantDTO `json:"variants,omitempty"`
}

// VariantDTO represents a product variant
type VariantDTO struct {
	BaseDTO
	ProductID     uint                  `json:"product_id"`
	SKU           string                `json:"sku"`
	Price         float64               `json:"price"`
	ComparePrice  float64               `json:"compare_price,omitempty"`
	StockQuantity int                   `json:"stock_quantity"`
	Attributes    []VariantAttributeDTO `json:"attributes"`
	Images        []string              `json:"images,omitempty"`
	IsDefault     bool                  `json:"is_default"`
}

type VariantAttributeDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateProductRequest represents the data needed to create a new product
type CreateProductRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Description   string                 `json:"description" validate:"required"`
	SKU           string                 `json:"sku" validate:"required"`
	Price         float64                `json:"price" validate:"required,gt=0"`
	StockQuantity int                    `json:"stock_quantity" validate:"required,gte=0"`
	Weight        float64                `json:"weight" validate:"required,gte=0"`
	CategoryID    uint                   `json:"category_id" validate:"required"`
	Images        []string               `json:"images"`
	Variants      []CreateVariantRequest `json:"variants,omitempty"`
}

// CreateVariantRequest represents the data needed to create a new product variant
type CreateVariantRequest struct {
	SKU           string                `json:"sku" validate:"required"`
	Price         float64               `json:"price" validate:"required,gt=0"`
	ComparePrice  float64               `json:"compare_price,omitempty"`
	StockQuantity int                   `json:"stock_quantity" validate:"required,gte=0"`
	Attributes    []VariantAttributeDTO `json:"attributes" validate:"required"`
	Images        []string              `json:"images,omitempty"`
	IsDefault     bool                  `json:"is_default"`
}

// UpdateProductRequest represents the data needed to update an existing product
type UpdateProductRequest struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	StockQuantity *int     `json:"stock_quantity,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	CategoryID    *uint    `json:"category_id,omitempty"`
	Images        []string `json:"images,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	ListResponseDTO[ProductDTO]
}

// ProductSearchRequest represents the parameters for searching products
type ProductSearchRequest struct {
	Query      string   `json:"query"`
	CategoryID *uint    `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	PaginationDTO
}
