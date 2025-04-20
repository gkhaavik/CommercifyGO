package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Define command line flags
	allFlag := flag.Bool("all", false, "Seed all data")
	usersFlag := flag.Bool("users", false, "Seed users data")
	categoriesFlag := flag.Bool("categories", false, "Seed categories data")
	productsFlag := flag.Bool("products", false, "Seed products data")
	ordersFlag := flag.Bool("orders", false, "Seed orders data")
	clearFlag := flag.Bool("clear", false, "Clear all data before seeding")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Clear data if requested
	if *clearFlag {
		if err := clearData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		fmt.Println("All data cleared")
	}

	// Seed data based on flags
	if *allFlag || *usersFlag {
		if err := seedUsers(db); err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		fmt.Println("Users seeded successfully")
	}

	if *allFlag || *categoriesFlag {
		if err := seedCategories(db); err != nil {
			log.Fatalf("Failed to seed categories: %v", err)
		}
		fmt.Println("Categories seeded successfully")
	}

	if *allFlag || *productsFlag {
		if err := seedProducts(db); err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}
		fmt.Println("Products seeded successfully")
	}

	if *allFlag || *ordersFlag {
		if err := seedOrders(db); err != nil {
			log.Fatalf("Failed to seed orders: %v", err)
		}
		fmt.Println("Orders seeded successfully")
	}

	if !*allFlag && !*usersFlag && !*categoriesFlag && !*productsFlag && !*ordersFlag && !*clearFlag {
		fmt.Println("No action specified")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
	}
}

// clearData clears all data from the database
func clearData(db *sql.DB) error {
	// Disable foreign key checks temporarily
	if _, err := db.Exec("SET CONSTRAINTS ALL DEFERRED"); err != nil {
		return err
	}

	// Clear tables in reverse order of dependencies
	tables := []string{
		"order_items",
		"orders",
		"cart_items",
		"carts",
		"products",
		"categories",
		"users",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return err
		}
		// Reset sequence
		if _, err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	if _, err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE"); err != nil {
		return err
	}

	return nil
}

// seedUsers seeds user data
func seedUsers(db *sql.DB) error {
	// Hash passwords
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userPassword, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	sellerPassword, err := bcrypt.GenerateFromPassword([]byte("seller123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()

	// Insert users
	users := []struct {
		email     string
		password  []byte
		firstName string
		lastName  string
		role      string
	}{
		{"admin@example.com", adminPassword, "Admin", "User", "admin"},
		{"user@example.com", userPassword, "Regular", "User", "user"},
		{"seller@example.com", sellerPassword, "Seller", "User", "seller"},
	}

	for _, user := range users {
		_, err := db.Exec(
			`INSERT INTO users (email, password, first_name, last_name, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (email) DO NOTHING`,
			user.email, user.password, user.firstName, user.lastName, user.role, now, now,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedCategories seeds category data
func seedCategories(db *sql.DB) error {
	now := time.Now()

	// Insert parent categories
	parentCategories := []struct {
		name        string
		description string
	}{
		{"Electronics", "Electronic devices and accessories"},
		{"Clothing", "Apparel and fashion items"},
		{"Home & Kitchen", "Home goods and kitchen appliances"},
		{"Books", "Books and publications"},
		{"Sports & Outdoors", "Sports equipment and outdoor gear"},
	}

	for _, category := range parentCategories {
		_, err := db.Exec(
			`INSERT INTO categories (name, description, parent_id, created_at, updated_at)
			VALUES ($1, $2, NULL, $3, $4)`,
			category.name, category.description, now, now,
		)
		if err != nil {
			return err
		}
	}

	// Get parent category IDs
	rows, err := db.Query("SELECT id, name FROM categories WHERE parent_id IS NULL")
	if err != nil {
		return err
	}
	defer rows.Close()

	parentCategoryIDs := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		parentCategoryIDs[name] = id
	}

	// Insert subcategories
	subcategories := []struct {
		name        string
		description string
		parentName  string
	}{
		{"Smartphones", "Mobile phones and accessories", "Electronics"},
		{"Laptops", "Notebook computers", "Electronics"},
		{"Audio", "Headphones, speakers, and audio equipment", "Electronics"},
		{"Men's Clothing", "Clothing for men", "Clothing"},
		{"Women's Clothing", "Clothing for women", "Clothing"},
		{"Footwear", "Shoes and boots", "Clothing"},
		{"Kitchen Appliances", "Appliances for the kitchen", "Home & Kitchen"},
		{"Furniture", "Home furniture", "Home & Kitchen"},
		{"Fiction", "Fiction books", "Books"},
		{"Non-Fiction", "Non-fiction books", "Books"},
		{"Fitness Equipment", "Equipment for exercise and fitness", "Sports & Outdoors"},
		{"Outdoor Gear", "Gear for outdoor activities", "Sports & Outdoors"},
	}

	for _, subcategory := range subcategories {
		parentID, ok := parentCategoryIDs[subcategory.parentName]
		if !ok {
			continue
		}

		_, err := db.Exec(
			`INSERT INTO categories (name, description, parent_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)`,
			subcategory.name, subcategory.description, parentID, now, now,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedProducts seeds product data
func seedProducts(db *sql.DB) error {
	// Get seller ID
	var sellerID int
	err := db.QueryRow("SELECT id FROM users WHERE role = 'seller' LIMIT 1").Scan(&sellerID)
	if err != nil {
		return err
	}

	// Get category IDs
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return err
	}
	defer rows.Close()

	categoryIDs := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		categoryIDs[name] = id
	}

	now := time.Now()

	// Insert products
	products := []struct {
		name         string
		description  string
		price        float64
		stock        int
		categoryName string
		images       string
	}{
		{
			"iPhone 13",
			"Apple iPhone 13 with A15 Bionic chip",
			999.99,
			50,
			"Smartphones",
			`["iphone13.jpg"]`,
		},
		{
			"Samsung Galaxy S21",
			"Samsung Galaxy S21 with 5G capability",
			899.99,
			75,
			"Smartphones",
			`["galaxys21.jpg"]`,
		},
		{
			"MacBook Pro",
			"Apple MacBook Pro with M1 chip",
			1299.99,
			30,
			"Laptops",
			`["macbookpro.jpg"]`,
		},
		{
			"Dell XPS 13",
			"Dell XPS 13 with Intel Core i7",
			1199.99,
			25,
			"Laptops",
			`["dellxps13.jpg"]`,
		},
		{
			"Sony WH-1000XM4",
			"Sony noise-cancelling headphones",
			349.99,
			100,
			"Audio",
			`["sonywh1000xm4.jpg"]`,
		},
		{
			"Men's Casual Shirt",
			"Comfortable casual shirt for men",
			39.99,
			200,
			"Men's Clothing",
			`["mencasualshirt.jpg"]`,
		},
		{
			"Women's Summer Dress",
			"Lightweight summer dress for women",
			49.99,
			150,
			"Women's Clothing",
			`["womendress.jpg"]`,
		},
		{
			"Running Shoes",
			"Comfortable shoes for running",
			89.99,
			120,
			"Footwear",
			`["runningshoes.jpg"]`,
		},
		{
			"Coffee Maker",
			"Automatic coffee maker for home use",
			79.99,
			80,
			"Kitchen Appliances",
			`["coffeemaker.jpg"]`,
		},
		{
			"Sofa Set",
			"3-piece sofa set for living room",
			599.99,
			15,
			"Furniture",
			`["sofaset.jpg"]`,
		},
		{
			"The Great Gatsby",
			"Classic novel by F. Scott Fitzgerald",
			12.99,
			300,
			"Fiction",
			`["greatgatsby.jpg"]`,
		},
		{
			"Atomic Habits",
			"Self-improvement book by James Clear",
			14.99,
			250,
			"Non-Fiction",
			`["atomichabits.jpg"]`,
		},
		{
			"Yoga Mat",
			"Non-slip yoga mat for exercise",
			24.99,
			180,
			"Fitness Equipment",
			`["yogamat.jpg"]`,
		},
		{
			"Camping Tent",
			"4-person camping tent for outdoor adventures",
			129.99,
			60,
			"Outdoor Gear",
			`["campingtent.jpg"]`,
		},
	}

	for _, product := range products {
		categoryID, ok := categoryIDs[product.categoryName]
		if !ok {
			continue
		}

		_, err := db.Exec(
			`INSERT INTO products (name, description, price, stock, category_id, seller_id, images, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			product.name, product.description, product.price, product.stock, categoryID, sellerID, product.images, now, now,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedOrders seeds order data
func seedOrders(db *sql.DB) error {
	// Get user IDs
	rows, err := db.Query("SELECT id FROM users WHERE role = 'user' OR role = 'admin'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("no users found to create orders for")
	}

	// Get product data
	productRows, err := db.Query("SELECT id, price FROM products")
	if err != nil {
		return err
	}
	defer productRows.Close()

	type productInfo struct {
		id    int
		price float64
	}

	var products []productInfo
	for productRows.Next() {
		var p productInfo
		if err := productRows.Scan(&p.id, &p.price); err != nil {
			return err
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found to create orders with")
	}

	// Sample addresses
	addresses := []map[string]string{
		{
			"street":      "123 Main St",
			"city":        "New York",
			"state":       "NY",
			"postal_code": "10001",
			"country":     "USA",
		},
		{
			"street":      "456 Oak Ave",
			"city":        "Los Angeles",
			"state":       "CA",
			"postal_code": "90001",
			"country":     "USA",
		},
		{
			"street":      "789 Pine Rd",
			"city":        "Chicago",
			"state":       "IL",
			"postal_code": "60601",
			"country":     "USA",
		},
		{
			"street":      "101 Maple Dr",
			"city":        "Seattle",
			"state":       "WA",
			"postal_code": "98101",
			"country":     "USA",
		},
	}

	// Order statuses
	statuses := []string{"pending", "paid", "shipped", "delivered", "cancelled"}

	// Payment providers
	paymentProviders := []string{"stripe", "paypal", "mock"}

	// Create orders
	for i := 0; i < 10; i++ {
		// Select random user
		userID := userIDs[i%len(userIDs)]

		// Select random address
		addrIndex := i % len(addresses)
		shippingAddr := addresses[addrIndex]
		billingAddr := addresses[addrIndex] // Use same address for billing

		// Convert addresses to JSON
		shippingAddrJSON, err := json.Marshal(shippingAddr)
		if err != nil {
			return err
		}

		billingAddrJSON, err := json.Marshal(billingAddr)
		if err != nil {
			return err
		}

		// Select random status
		status := statuses[i%len(statuses)]

		// Create timestamps
		now := time.Now()
		createdAt := now.Add(time.Duration(-i*24) * time.Hour) // Each order created a day apart
		updatedAt := createdAt

		// Set completed_at for delivered orders
		var completedAt *time.Time
		if status == "delivered" {
			completedTime := updatedAt.Add(3 * 24 * time.Hour) // 3 days after creation
			completedAt = &completedTime
		}

		// Set payment details for paid, shipped, or delivered orders
		var paymentID string
		var paymentProvider string
		var trackingCode string

		if status == "paid" || status == "shipped" || status == "delivered" {
			paymentID = fmt.Sprintf("payment_%d_%s", i, time.Now().Format("20060102"))
			paymentProvider = paymentProviders[i%len(paymentProviders)]
		}

		if status == "shipped" || status == "delivered" {
			trackingCode = fmt.Sprintf("TRACK%d%s", i, time.Now().Format("20060102"))
		}

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		// Insert order
		var orderID int
		err = tx.QueryRow(`
			INSERT INTO orders (
				user_id, total_amount, status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id
		`,
			userID, 0, // Total amount will be updated after adding items
			status,
			shippingAddrJSON,
			billingAddrJSON,
			paymentID,
			paymentProvider,
			trackingCode,
			createdAt,
			updatedAt,
			completedAt,
		).Scan(&orderID)

		if err != nil {
			tx.Rollback()
			return err
		}

		// Add 1-3 random products as order items
		numItems := (i % 3) + 1
		totalAmount := 0.0

		// Ensure we don't try to add more items than we have products
		if numItems > len(products) {
			numItems = len(products)
		}

		for j := 0; j < numItems; j++ {
			// Select product
			product := products[(i+j)%len(products)]

			// Random quantity between 1 and 3
			quantity := (j % 3) + 1

			// Calculate subtotal
			subtotal := float64(quantity) * product.price
			totalAmount += subtotal

			// Insert order item
			_, err = tx.Exec(`
				INSERT INTO order_items (
					order_id, product_id, quantity, price, subtotal, created_at
				)
				VALUES ($1, $2, $3, $4, $5, $6)
			`,
				orderID,
				product.id,
				quantity,
				product.price,
				subtotal,
				createdAt,
			)

			if err != nil {
				tx.Rollback()
				return err
			}
		}

		// Update order with total amount
		_, err = tx.Exec(`
			UPDATE orders
			SET total_amount = $1
			WHERE id = $2
		`,
			totalAmount,
			orderID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
