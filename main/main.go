package main

import (
	"Naol_paxosLab/paxos"
	"fmt"
)

func main() {
	acceptors := []*paxos.Acceptor{
		&paxos.Acceptor{},
		&paxos.Acceptor{},
		&paxos.Acceptor{},
	}
	proposer := paxos.Proposer{ProposalNumber: 1, Value: "Distributed Systems"}
	value := proposer.Propose("Distributed Systems", acceptors)
	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}
}
