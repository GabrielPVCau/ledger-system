package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const BaseURL = "http://localhost:8080"

type TransferRequest struct {
	From   int   `json:"from"`
	To     int   `json:"to"`
	Amount int64 `json:"amount"` // using int64 for correct json marshaling
}

func main() {
	log.Println("Starting E2E Functional Tests (Edge Cases)...")
	failed := false

	// Scenario A: Negative Amount
	if !runTest("Scenario A: Negative Amount", TransferRequest{From: 1, To: 2, Amount: -100}, http.StatusBadRequest, http.StatusUnprocessableEntity) {
		failed = true
	}

	// Scenario B: Insufficient Funds
	// Assuming Alice (1) has 10000. Try to send 20000.
	if !runTest("Scenario B: Insufficient Funds", TransferRequest{From: 1, To: 2, Amount: 9000000}, http.StatusBadRequest, http.StatusUnprocessableEntity) {
		failed = true
	}

	// Scenario C: Non-existent Account
	if !runTest("Scenario C: Non-existent Account", TransferRequest{From: 1, To: 99999, Amount: 100}, http.StatusNotFound, http.StatusBadRequest) {
		failed = true
	}

	// Scenario D: Valid Transfer
	if !runTest("Scenario D: Valid Transfer", TransferRequest{From: 1, To: 2, Amount: 10}, http.StatusOK, http.StatusCreated) {
		failed = true
	}

	if failed {
		log.Println("❌ Some tests FAILED. Fix the code!")
		os.Exit(1)
	} else {
		log.Println("✅ All E2E tests PASSED!")
	}
}

func runTest(name string, payload TransferRequest, expectedStatus ...int) bool {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(BaseURL+"/transfer", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ %s: Failed to make request: %v", name, err)
		return false
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	for _, s := range expectedStatus {
		if resp.StatusCode == s {
			log.Printf("✅ %s: PASSED (Status: %d)", name, resp.StatusCode)
			return true
		}
	}

	log.Printf("❌ %s: FAILED. Expected one of %v, got %d. Body: %s", name, expectedStatus, resp.StatusCode, string(respBody))
	return false
}
