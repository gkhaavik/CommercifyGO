package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// ShippingHandler handles shipping-related HTTP requests
type ShippingHandler struct {
	shippingUseCase *usecase.ShippingUseCase
	logger          logger.Logger
}

// NewShippingHandler creates a new ShippingHandler
func NewShippingHandler(shippingUseCase *usecase.ShippingUseCase, logger logger.Logger) *ShippingHandler {
	return &ShippingHandler{
		shippingUseCase: shippingUseCase,
		logger:          logger,
	}
}

// CalculateShippingOptions handles calculating available shipping options for an address and order details
func (h *ShippingHandler) CalculateShippingOptions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		Address     entity.Address `json:"address"`
		OrderValue  float64        `json:"order_value"`
		OrderWeight float64        `json:"order_weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate shipping options
	options, err := h.shippingUseCase.CalculateShippingOptions(
		requestBody.Address,
		requestBody.OrderValue,
		requestBody.OrderWeight,
	)
	if err != nil {
		h.logger.Error("Failed to calculate shipping options: %v", err)
		http.Error(w, "Failed to calculate shipping options", http.StatusInternalServerError)
		return
	}

	// Return shipping options
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// GetShippingMethodByID handles retrieving a shipping method by ID
func (h *ShippingHandler) GetShippingMethodByID(w http.ResponseWriter, r *http.Request) {
	// Get method ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping method ID", http.StatusBadRequest)
		return
	}

	// Get shipping method
	method, err := h.shippingUseCase.GetShippingMethodByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping method: %v", err)
		http.Error(w, "Shipping method not found", http.StatusNotFound)
		return
	}

	// Return shipping method
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(method)
}

// ListShippingMethods handles listing all shipping methods
func (h *ShippingHandler) ListShippingMethods(w http.ResponseWriter, r *http.Request) {
	// Get active parameter from query string
	activeOnly := r.URL.Query().Get("active") == "true"

	// Get shipping methods
	methods, err := h.shippingUseCase.ListShippingMethods(activeOnly)
	if err != nil {
		h.logger.Error("Failed to list shipping methods: %v", err)
		http.Error(w, "Failed to list shipping methods", http.StatusInternalServerError)
		return
	}

	// Return shipping methods
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}

// CreateShippingMethod handles creating a new shipping method (admin only)
func (h *ShippingHandler) CreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateShippingMethodInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create shipping method
	method, err := h.shippingUseCase.CreateShippingMethod(input)
	if err != nil {
		h.logger.Error("Failed to create shipping method: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created shipping method
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(method)
}

// UpdateShippingMethod handles updating a shipping method (admin only)
func (h *ShippingHandler) UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Get method ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping method ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input usecase.UpdateShippingMethodInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID from URL
	input.ID = uint(id)

	// Update shipping method
	method, err := h.shippingUseCase.UpdateShippingMethod(input)
	if err != nil {
		h.logger.Error("Failed to update shipping method: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated shipping method
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(method)
}

// CreateShippingZone handles creating a new shipping zone (admin only)
func (h *ShippingHandler) CreateShippingZone(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateShippingZoneInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create shipping zone
	zone, err := h.shippingUseCase.CreateShippingZone(input)
	if err != nil {
		h.logger.Error("Failed to create shipping zone: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created shipping zone
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(zone)
}

// GetShippingZoneByID handles retrieving a shipping zone by ID
func (h *ShippingHandler) GetShippingZoneByID(w http.ResponseWriter, r *http.Request) {
	// Get zone ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping zone ID", http.StatusBadRequest)
		return
	}

	// Get shipping zone
	zone, err := h.shippingUseCase.GetShippingZoneByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping zone: %v", err)
		http.Error(w, "Shipping zone not found", http.StatusNotFound)
		return
	}

	// Return shipping zone
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

// ListShippingZones handles listing all shipping zones
func (h *ShippingHandler) ListShippingZones(w http.ResponseWriter, r *http.Request) {
	// Get active parameter from query string
	activeOnly := r.URL.Query().Get("active") == "true"

	// Get shipping zones
	zones, err := h.shippingUseCase.ListShippingZones(activeOnly)
	if err != nil {
		h.logger.Error("Failed to list shipping zones: %v", err)
		http.Error(w, "Failed to list shipping zones", http.StatusInternalServerError)
		return
	}

	// Return shipping zones
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

// UpdateShippingZone handles updating a shipping zone (admin only)
func (h *ShippingHandler) UpdateShippingZone(w http.ResponseWriter, r *http.Request) {
	// Get zone ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping zone ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input usecase.UpdateShippingZoneInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID from URL
	input.ID = uint(id)

	// Update shipping zone
	zone, err := h.shippingUseCase.UpdateShippingZone(input)
	if err != nil {
		h.logger.Error("Failed to update shipping zone: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated shipping zone
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

// CreateShippingRate handles creating a new shipping rate (admin only)
func (h *ShippingHandler) CreateShippingRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateShippingRateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create shipping rate
	rate, err := h.shippingUseCase.CreateShippingRate(input)
	if err != nil {
		h.logger.Error("Failed to create shipping rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created shipping rate
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rate)
}

// GetShippingRateByID handles retrieving a shipping rate by ID
func (h *ShippingHandler) GetShippingRateByID(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Get shipping rate
	rate, err := h.shippingUseCase.GetShippingRateByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping rate: %v", err)
		http.Error(w, "Shipping rate not found", http.StatusNotFound)
		return
	}

	// Return shipping rate
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

// UpdateShippingRate handles updating a shipping rate (admin only)
func (h *ShippingHandler) UpdateShippingRate(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input usecase.UpdateShippingRateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set ID from URL
	input.ID = uint(id)

	// Update shipping rate
	rate, err := h.shippingUseCase.UpdateShippingRate(input)
	if err != nil {
		h.logger.Error("Failed to update shipping rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated shipping rate
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

// CreateWeightBasedRate handles creating a new weight-based shipping rate (admin only)
func (h *ShippingHandler) CreateWeightBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateWeightBasedRateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create weight-based rate
	rate, err := h.shippingUseCase.CreateWeightBasedRate(input)
	if err != nil {
		h.logger.Error("Failed to create weight-based rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created weight-based rate
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rate)
}

// CreateValueBasedRate handles creating a new value-based shipping rate (admin only)
func (h *ShippingHandler) CreateValueBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input usecase.CreateValueBasedRateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create value-based rate
	rate, err := h.shippingUseCase.CreateValueBasedRate(input)
	if err != nil {
		h.logger.Error("Failed to create value-based rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created value-based rate
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rate)
}

// GetShippingCost handles calculating shipping cost for a specific shipping rate
func (h *ShippingHandler) GetShippingCost(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var requestBody struct {
		OrderValue  float64 `json:"order_value"`
		OrderWeight float64 `json:"order_weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate shipping cost
	cost, err := h.shippingUseCase.GetShippingCost(
		uint(id),
		requestBody.OrderValue,
		requestBody.OrderWeight,
	)
	if err != nil {
		h.logger.Error("Failed to calculate shipping cost: %v", err)
		http.Error(w, "Failed to calculate shipping cost", http.StatusInternalServerError)
		return
	}

	// Return shipping cost
	response := map[string]float64{
		"cost": cost,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
