package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type RazorpayClient struct {
	KeyID     string
	KeySecret string
}

func NewRazorpayClient() *RazorpayClient {
	return &RazorpayClient{
		KeyID:     os.Getenv("RAZORPAY_KEY_ID"),
		KeySecret: os.Getenv("RAZORPAY_KEY_SECRET"),
	}
}

type CreateOrderRequest struct {
	Amount         int64  `json:"amount"` 
	Currency       string `json:"currency"`
	Receipt        string `json:"receipt"`
	PaymentCapture int    `json:"payment_capture"` 
}

type CreateOrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (c *RazorpayClient) CreateOrder(amount float64, receiptID string) (*CreateOrderResponse, error) {
	reqBody := CreateOrderRequest{
		Amount:         int64(amount * 100), 
		Currency:       "INR",
		Receipt:        receiptID,
		PaymentCapture: 1,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.KeyID, c.KeySecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New("razorpay order creation failed: " + string(respBody))
	}

	var orderRes CreateOrderResponse
	err = json.Unmarshal(respBody, &orderRes)
	if err != nil {
		return nil, err
	}
	return &orderRes, nil
}
