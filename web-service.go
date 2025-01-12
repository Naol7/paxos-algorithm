package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"Naol_paxosLab/paxos"
)

var (
	acceptors = []*paxos.Acceptor{{}, {}, {}}
	mu        sync.Mutex
)

// proposeWithRetry handles retries and timeouts for Paxos proposals.
func proposeWithRetry(ctx context.Context, proposer paxos.Proposer, value interface{}, acceptors []*paxos.Acceptor, retries int) interface{} {
	for i := 0; i <= retries; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Propose timed out after %d retries\n", i)
			return nil
		default:
			mu.Lock()
			result := proposer.Propose(value, acceptors)
			mu.Unlock()

			if result != nil {
				return result
			}

			log.Printf("Retrying proposal (attempt %d)\n", i+1)
			time.Sleep(500 * time.Millisecond) // Retry delay
		}
	}
	log.Println("Propose failed after maximum retries")
	return nil
}

func proposeHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProposalNumber int
		Value          string
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	proposer := paxos.Proposer{ProposalNumber: body.ProposalNumber, Value: body.Value}
	value := proposeWithRetry(ctx, proposer, body.Value, acceptors, 3)

	if value != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %s\n", value)
	} else {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Consensus not reached\n")
	}
}

func main() {
	http.HandleFunc("/propose", proposeHandler)
	log.Println("Paxos web service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
