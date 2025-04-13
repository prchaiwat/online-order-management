package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"order-management-system/internal/config"
	"order-management-system/internal/model"
	"order-management-system/internal/repository"
	"order-management-system/internal/service"
	"order-management-system/internal/util"
)

type OrderHandler struct {
	service service.OrderService
	config  *config.Config
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return repository.NewOrderRepository(db)
}

func NewOrderService(repo repository.OrderRepository, config *config.Config) service.OrderService {
	return service.NewOrderService(repo, config)
}

func NewOrderHandler(service service.OrderService, config *config.Config) *OrderHandler {
	return &OrderHandler{
		service: service,
		config:  config,
	}
}

func (h *OrderHandler) CreateOrderWithOutGoroutine(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.RequestTimeoutSec)*time.Second)
	defer cancel()

	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		handleError(w, "json format incorrect", http.StatusBadRequest)
		return
	}
	orderJson, _ := json.Marshal(order)
	log.Println("/request CreateOrderWithOutGoroutine body:", string(orderJson))

	err = util.ValidateOrder(order)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := h.service.CreateOrderWithOutGoroutine(ctx, order)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"order_id": orderID,
		"status":   "created",
	}
	responseJson, _ := json.Marshal(response)
	log.Println("response CreateOrderWithOutGoroutine :", string(responseJson))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		handleError(w, "json format incorrect", http.StatusBadRequest)
		return
	}

	orderJson, _ := json.Marshal(order)
	log.Println("/request CreateOrder body:", string(orderJson))

	err = util.ValidateOrder(order)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := h.service.CreateOrder(order)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"order_id": orderID,
		"status":   "created",
	}
	responseJson, _ := json.Marshal(response)
	log.Println("response :", string(responseJson))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.RequestTimeoutSec)*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		handleError(w, "order ID format incorrect", http.StatusBadRequest)
		return
	}

	log.Println("/request GetOrder parameter:", orderID)

	order, err := h.service.GetOrder(ctx, orderID)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.OrderByIdResponse{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		TotalAmount:  order.TotalAmount,
		Items:        order.Items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.RequestTimeoutSec)*time.Second)
	defer cancel()

	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	// default values
	page := 1
	size := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr != "" {
		s, err := strconv.Atoi(sizeStr)
		if err == nil && s > 0 {
			size = s
		}
	}
	log.Printf("/request GetOrders parameter page: %d ,size: %d", page, size)

	response, err := h.service.GetOrders(ctx, page, size)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.RequestTimeoutSec)*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		handleError(w, "order ID format incorrect", http.StatusBadRequest)
		return
	}

	log.Println("/request UpdateOrderStatus parameter:", orderID)

	var orderStatus model.OrderStatus
	err = json.NewDecoder(r.Body).Decode(&orderStatus)
	if err != nil {
		handleError(w, "Json format incorrect", http.StatusBadRequest)
		return
	}

	if orderStatus.Status == "" {
		handleError(w, "status can't be empty", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrderStatus(ctx, orderID, orderStatus.Status)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.OrderStatusResponse{
		OrderID: orderID,
		Status:  orderStatus.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleError(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	response := struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	}{
		StatusCode: statusCode,
		Message:    err,
	}

	json.NewEncoder(w).Encode(response)
}
