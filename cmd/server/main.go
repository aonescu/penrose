package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aonescu/penrose/encryption"

	"github.com/aonescu/penrose/utils"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
)

func main() {
	// Parse flags
	now := time.Now()
	cc := flag.String("cc", "", "Path to context configuration file")
	evkFile := flag.String("key_eval", "", "Path to evaluation key file")
	inputFile := flag.String("input", "", "Path to input file")
	outputFile := flag.String("output", "", "Path to output file")

	flag.Parse()

	// Initialize required variables
	params := utils.Parameters{}
	evk := utils.EvaluationKeySet{}
	in := rlwe.NewCiphertext(params, 1, params.MaxLevel)

	// Deserialize configuration, keys, and input ciphertext
	if err := utils.Deserialize(&params, *cc); err != nil {
		log.Fatalf("Error deserializing params: %v", err)
	}

	if err := utils.Deserialize(&evk, *evkFile); err != nil {
		log.Fatalf("Error deserializing evaluation key: %v", err)
	}

	if err := utils.Deserialize(&in, *inputFile); err != nil {
		log.Fatalf("Error deserializing input ciphertext: %v", err)
	}

	// Initialize Lattigo encryption context from configuration
	lattigoContext, err := encryption.NewLattigoContext(*cc)
	if err != nil {
		log.Fatalf("Error initializing Lattigo context: %v", err)
	}

	// Example of using the encryption context to perform homomorphic operations
	// Assuming that the solution requires some homomorphic operations
	// like adding or multiplying encrypted bids
	// Here, just a simple decrypt and re-encrypt to show the flow

	// Decrypt the input bid value
	decryptedBid := lattigoContext.DecryptBid(in)
	fmt.Printf("Decrypted Bid: %f\n", decryptedBid)

	// You can perform other operations like adding encrypted bids, etc.
	// encryptedBid := lattigoContext.EncryptBid(100.0)  // example of encrypting a bid
	// encryptedResult := lattigoContext.AddBids(&in, encryptedBid)

	// Call the solution to process the decrypted or encrypted data (based on actual use case)
	solution := &Solution{}
	out, err := solution.SolveTestcase(params, evk, in) // Pass in directly if it's *rlwe.Ciphertext
	// Process the decrypted or encrypted data (based on actual use case)
	// out := processTestcase(params, evk, &in)
	// if err != nil {
	// 	log.Fatalf("processTestcase error: %v", err)
	// }
	// Serialize the output result
	if err := utils.Serialize(out, *outputFile); err != nil {
		log.Fatalf("Error serializing output: %v", err)
	}

	// Serialize the output result
	if err := utils.Serialize(out, *outputFile); err != nil {
		log.Fatalf("Error serializing output: %v", err)
	}

	// Print execution time
	fmt.Printf("Done: %s\n", time.Since(now))
}

type Solution struct{}

func (s *Solution) SolveTestcase(params utils.Parameters, evk utils.EvaluationKeySet, in *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	// Implement the solution logic here
	// For example, perform some homomorphic operations on the input ciphertext
	// and return the resulting ciphertext

	// Placeholder logic: just return the input ciphertext as output
	return in, nil
}
