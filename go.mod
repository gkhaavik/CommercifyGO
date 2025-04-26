module github.com/zenfulcode/commercify

go 1.24.0

require (
	github.com/gkhaavik/vipps-mobilepay-sdk v0.0.2
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.10.0
	github.com/stripe/stripe-go/v72 v72.122.0
	golang.org/x/crypto v0.37.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gkhaavik/vipps-mobilepay-sdk => ../vipps-mobilepay-sdk
