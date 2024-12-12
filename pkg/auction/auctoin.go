package auction

import (
	"github.com/aonescu/penrose/encryption"

	"github.com/tuneinsight/lattigo/ckks"
)

type Auction struct {
	EncryptedBids map[string]*ckks.Ciphertext // Map of bidder ID to encrypted bid
	Context       *encryption.LattigoContext
}

// NewAuction creates a new auction instance.
func NewAuction(ctx *encryption.LattigoContext) *Auction {
	return &Auction{
		EncryptedBids: make(map[string]*ckks.Ciphertext),
		Context:       ctx,
	}
}

func (a *Auction) DetermineWinner() string {
	var highestBid float64
	var winner string

	for bidderID, encryptedBid := range a.EncryptedBids {
		// Decrypt the ciphertext using the decryptor from the ckks package
		plaintext := a.Context.Decryptor.DecryptNew(encryptedBid)
		// Decode the plaintext using the ckks decoder
		decodedValue := real(a.Context.Encoder.Decode(plaintext)[0]) // Example: decode real part
		if decodedValue > highestBid {
			highestBid = decodedValue
			winner = bidderID
		}
	}

	return winner
}