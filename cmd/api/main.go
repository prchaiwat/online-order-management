package main

import (
	"log"
	"net/http"
	"order-management-system/internal/config"
	"order-management-system/internal/handler"
	"order-management-system/internal/middleware"
	"order-management-system/pkg/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("can't load env:", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("can't load config:", err)
		return
	}
	database, err := db.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("can't connect database:", err)
		return
	}
	defer database.Close()
	log.Println("connect database postgres success")

	orderRepo := handler.NewOrderRepository(database)
	orderService := handler.NewOrderService(orderRepo, cfg)
	orderHandler := handler.NewOrderHandler(orderService, cfg)

	router := mux.NewRouter()

	router.Use(middleware.CorsMiddleware)

	router.HandleFunc("/ordersWithOutGoroutine", orderHandler.CreateOrderWithOutGoroutine).Methods("POST")
	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("OPTIONS")
	router.HandleFunc("/orders/{order_id}", orderHandler.GetOrder).Methods("GET")
	router.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{order_id}/status", orderHandler.UpdateOrderStatus).Methods("PUT")
	router.HandleFunc("/orders/{order_id}/status", orderHandler.UpdateOrderStatus).Methods("OPTIONS")

	log.Println("server start at port:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
