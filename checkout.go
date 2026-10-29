package lyngbak

import (
	"encoding/json"
	"fmt"
)

type CheckoutRequiredError struct {
	Message     string
	CheckoutURL string
}

func (e *CheckoutRequiredError) Error() string {
	return fmt.Sprintf("lyngbak: checkout required: %s (%s)", e.Message, e.CheckoutURL)
}

type PaymentPendingError struct {
	Message string
	OrderID string
}

func (e *PaymentPendingError) Error() string {
	return fmt.Sprintf("lyngbak: payment pending: %s (order %s)", e.Message, e.OrderID)
}

func decodeResourceOrPayment(raw json.RawMessage) (*Resource, error) {
	if len(raw) == 0 {
		return new(Resource), nil
	}
	var probe struct {
		UUID        string `json:"uuid"`
		Message     string `json:"message"`
		CheckoutURL string `json:"checkout_url"`
		OrderID     string `json:"order_id"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, fmt.Errorf("lyngbak: decoding response: %w", err)
	}
	if probe.UUID == "" {
		if probe.CheckoutURL != "" {
			return nil, &CheckoutRequiredError{Message: probe.Message, CheckoutURL: probe.CheckoutURL}
		}
		if probe.OrderID != "" {
			return nil, &PaymentPendingError{Message: probe.Message, OrderID: probe.OrderID}
		}
	}
	out := new(Resource)
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, fmt.Errorf("lyngbak: decoding response: %w", err)
	}
	return out, nil
}
