package encryption

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aonescu/penrose/utils"
	"github.com/tuneinsight/lattigo/v5/he/hefloat/bootstrapping"
	"github.com/tuneinsight/lattigo/v5/rlwe"
)

// LattigoContext encapsulates all Lattigo-related operations and parameters
type LattigoContext struct {
	Scheme        *rlwe.Parameters
	Bootstrapping *bootstrapping.Parameters
	KeyGen        *rlwe.KeyGenerator
	Encryptor     *rlwe.Encryptor
	Decryptor     *rlwe.Decryptor
	Evaluator     *rlwe.Evaluator
}

// NewLattigoContext initializes a new Lattigo context based on a configuration file
func NewLattigoContext(configFilePath string) (*LattigoContext, error) {
	// Load the configuration file
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config struct {
		Scheme struct {
			LogN            int   `json:"LogN"`
			LogQ            []int `json:"LogQ"`
			LogP            []int `json:"LogP"`
			LogDefaultScale int   `json:"LogDefaultScale"`
		} `json:"Scheme"`
		Bootstrapping struct {
			Enable bool `json:"Enable"`
		} `json:"Bootstrapping"`
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Create the parameters using the config
	scheme := rlwe.NewParameters(config.Scheme.LogN, config.Scheme.LogQ, config.Scheme.LogP, config.Scheme.LogDefaultScale)

	// Initialize the Bootstrapping parameters
	var bootParams *bootstrapping.Parameters
	if config.Bootstrapping.Enable {
		bootParams = bootstrapping.NewParameters(scheme)
	}

	// Create the KeyGenerator, Encryptor, Decryptor, and Evaluator
	keyGen := rlwe.NewKeyGenerator(scheme)
	sk, pk := keyGen.GenKeyPair()
	encryptor := rlwe.NewEncryptor(scheme, pk)
	decryptor := rlwe.NewDecryptor(scheme, sk)
	evaluator := rlwe.NewEvaluator(scheme)

	return &LattigoContext{
		Scheme:        scheme,
		Bootstrapping: bootParams,
		KeyGen:        keyGen,
		Encryptor:     encryptor,
		Decryptor:     decryptor,
		Evaluator:     evaluator,
	}, nil
}

// Serialize writes the LattigoContext parameters to a file
func (ctx *LattigoContext) Serialize(path string) error {
	return utils.Serialize(ctx, path)
}

// DeserializeLattigoContext loads the LattigoContext from a file
func DeserializeLattigoContext(path string) (*LattigoContext, error) {
	var ctx LattigoContext
	if err := utils.Deserialize(&ctx, path); err != nil {
		return nil, fmt.Errorf("failed to deserialize LattigoContext: %w", err)
	}
	return &ctx, nil
}

// EncryptBid encrypts a bid value
func (ctx *LattigoContext) EncryptBid(bidValue float64) *rlwe.Ciphertext {
	plaintext := rlwe.NewPlaintext(ctx.Scheme)
	ctx.Scheme.SlotEncoder().Encode([]float64{bidValue}, plaintext)
	return ctx.Encryptor.EncryptNew(plaintext)
}

// DecryptBid decrypts an encrypted bid value
func (ctx *LattigoContext) DecryptBid(ciphertext *rlwe.Ciphertext) float64 {
	plaintext := ctx.Decryptor.DecryptNew(ciphertext)
	values := ctx.Scheme.SlotEncoder().Decode(plaintext, ctx.Scheme.LogSlots())
	return values[0]
}

// AddBids adds two encrypted bid values
func (ctx *LattigoContext) AddBids(c1, c2 *rlwe.Ciphertext) *rlwe.Ciphertext {
	return ctx.Evaluator.AddNew(c1, c2)
}
