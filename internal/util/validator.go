package util

import (
	"fmt"
	"strings"

	"order-management-system/internal/model"
)

func ValidateOrder(order model.Order) error {

	if len(strings.TrimSpace(order.CustomerName)) == 0 {
		return fmt.Errorf("have to input customer name")
	}

	if len(order.CustomerName) > 100 {
		return fmt.Errorf("customer name must not exceed 100 chars")
	}

	if len(order.Items) == 0 {
		return fmt.Errorf("have to input at least one item")
	}

	var totalAmount float64
	for i, item := range order.Items {
		if len(strings.TrimSpace(item.ProductName)) == 0 {
			return fmt.Errorf("order item %d: have to input product name", i+1)
		}

		if len(item.ProductName) > 100 {
			return fmt.Errorf("order item %d: product name must not exceed 100 chars", i+1)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("order item %d: have to input quantity > 0", i+1)
		}

		if item.Price < 0 {
			return fmt.Errorf("order item %d: have to input price >= 0", i+1)
		}

		totalAmount += item.Price * float64(item.Quantity)
	}

	if totalAmount > 99999999.99 {
		return fmt.Errorf("Total price cannot exceed 99,999,999.99")
	}

	return nil
}
