package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Config
const (
	BaseURL        = "http://localhost:8080"
	AccountA       = 1
	AccountB       = 2
	Workers        = 50
	TransfersPerWorker = 20
	Amount         = 10 // 0.10 currency units
)

type TransferRequest struct {
	From   int   `json:"from"`
	To     int   `json:"to"`
	Amount int64 `json:"amount"`
}

func main() {
	log.Println("Starting stress test...")
	log.Printf("Workers: %d, Transfers/Worker: %d\n", Workers, TransfersPerWorker)

	var wg sync.WaitGroup
	wg.Add(Workers)

	start := time.Now()

	for i := 0; i < Workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < TransfersPerWorker; j++ {
				direction := j % 2 // Alternate direction
				from, to := AccountA, AccountB
				if direction == 1 {
					from, to = AccountB, AccountA
				}

				err := makeTransfer(from, to, Amount)
				if err != nil {
					log.Printf("Worker %d failed transfer: %v", workerID, err)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	
	log.Printf("Stress test completed in %v", duration)
	log.Println("Verify the database balances manually. Total balance of Account A + Account B should be unchanged from start (10000 + 0 = 10000).")
}

func makeTransfer(from, to int, amount int64) error {
	reqBody, _ := json.Marshal(TransferRequest{
		From:   from,
		To:     to,
		Amount: amount,
	})

	resp, err := http.Post(BaseURL+"/transfer", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return nil
}
