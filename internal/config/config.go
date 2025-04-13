package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL         string
	Port                string
	RequestTimeoutSec   int
	MaxConcurrentOrders int
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	timeoutStr := os.Getenv("REQUEST_TIMEOUT_SEC")
	timeout := 5 //default timeout
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 {
			timeout = t
		}
	}

	maxConcurrentOrdersStr := os.Getenv("MAX_CONCURRENT_ORDERS")
	maxConcurrentOrders := 20 //default max concurrent orders
	if maxConcurrentOrdersStr != "" {
		if m, err := strconv.Atoi(maxConcurrentOrdersStr); err == nil && m > 0 {
			maxConcurrentOrders = m
		}
	}

	return &Config{
		DatabaseURL:         dbURL,
		Port:                port,
		RequestTimeoutSec:   timeout,
		MaxConcurrentOrders: maxConcurrentOrders,
	}, nil
}
