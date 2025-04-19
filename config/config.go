package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Payment  PaymentConfig
	Email    EmailConfig
	Stripe   StripeConfig
	PayPal   PayPalConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AuthConfig holds authentication-specific configuration
type AuthConfig struct {
	JWTSecret     string
	TokenDuration int
}

// PaymentConfig holds payment-specific configuration
type PaymentConfig struct {
	EnabledProviders []string // List of enabled payment providers
}

// EmailConfig holds email-specific configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	AdminEmail   string
	Enabled      bool
}

// StripeConfig holds Stripe-specific configuration
type StripeConfig struct {
	SecretKey          string
	PublicKey          string
	WebhookSecret      string
	PaymentDescription string
	Enabled            bool
}

// PayPalConfig holds PayPal-specific configuration
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	Sandbox      bool
	Enabled      bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	readTimeout, err := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "15"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
	}

	tokenDuration, err := strconv.Atoi(getEnv("AUTH_TOKEN_DURATION", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_TOKEN_DURATION: %w", err)
	}

	smtpPort, err := strconv.Atoi(getEnv("EMAIL_SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid EMAIL_SMTP_PORT: %w", err)
	}

	emailEnabled, err := strconv.ParseBool(getEnv("EMAIL_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid EMAIL_ENABLED: %w", err)
	}

	stripeEnabled, err := strconv.ParseBool(getEnv("STRIPE_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid STRIPE_ENABLED: %w", err)
	}

	paypalEnabled, err := strconv.ParseBool(getEnv("PAYPAL_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid PAYPAL_ENABLED: %w", err)
	}

	paypalSandbox, err := strconv.ParseBool(getEnv("PAYPAL_SANDBOX", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid PAYPAL_SANDBOX: %w", err)
	}

	// Parse enabled payment providers
	enabledProviders := []string{"mock"} // Always enable mock provider for testing
	if stripeEnabled {
		enabledProviders = append(enabledProviders, "stripe")
	}
	if paypalEnabled {
		enabledProviders = append(enabledProviders, "paypal")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("AUTH_JWT_SECRET", "your-secret-key"),
			TokenDuration: tokenDuration,
		},
		Payment: PaymentConfig{
			EnabledProviders: enabledProviders,
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
			SMTPPort:     smtpPort,
			SMTPUsername: getEnv("EMAIL_SMTP_USERNAME", ""),
			SMTPPassword: getEnv("EMAIL_SMTP_PASSWORD", ""),
			FromEmail:    getEnv("EMAIL_FROM_ADDRESS", "noreply@example.com"),
			FromName:     getEnv("EMAIL_FROM_NAME", "E-Commerce Store"),
			AdminEmail:   getEnv("EMAIL_ADMIN_ADDRESS", "admin@example.com"),
			Enabled:      emailEnabled,
		},
		Stripe: StripeConfig{
			SecretKey:          getEnv("STRIPE_SECRET_KEY", ""),
			PublicKey:          getEnv("STRIPE_PUBLIC_KEY", ""),
			WebhookSecret:      getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PaymentDescription: getEnv("STRIPE_PAYMENT_DESCRIPTION", "E-Commerce Store Purchase"),
			Enabled:            stripeEnabled,
		},
		PayPal: PayPalConfig{
			ClientID:     getEnv("PAYPAL_CLIENT_ID", ""),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", ""),
			Sandbox:      paypalSandbox,
			Enabled:      paypalEnabled,
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
