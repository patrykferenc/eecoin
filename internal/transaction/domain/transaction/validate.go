package transaction

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

func validateCoinbase(tx *Transaction, blockHeight int) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil, the first transaction in a block must be a coinbase transaction")
	}
	if ins := len(tx.In); ins != 1 {
		return fmt.Errorf("coinbase transaction must have one input, got %d", ins)
	}
	if tx.In[0].OutputIndex() != blockHeight {
		return fmt.Errorf("coinbase transaction input must have output index equal to block height, got %d", tx.In[0].OutputIndex())
	}
	if outs := len(tx.Ou); outs != 1 {
		return fmt.Errorf("coinbase transaction must have one output, got %d", outs)
	}
	if tx.Ou[0].Amoun != COINBASE_AMOUNT {
		return fmt.Errorf("coinbase transaction output amount must be %d, got %d", COINBASE_AMOUNT, tx.Ou[0].Amoun)
	}

	return nil
}

func validateTransactionIn(inputTx *Input, tx *Transaction, unspent UnspentOutputRepository) error {
	referencedOutput, err := unspent.GetByOutputIDAndIndex(
		inputTx.OutputId,
		inputTx.OutputIdx,
	)
	if err != nil {
		return fmt.Errorf("error getting referenced output: %w", err)
	}
	if referencedOutput == (UnspentOutput{}) {
		return fmt.Errorf("referenced output not found")
	}

	// Decode the address (public key) from the referenced output
	pubKeyBytes, err := hex.DecodeString(referencedOutput.Address())
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Parse the public key
	publicKeyRaw, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	publicKey, ok := publicKeyRaw.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to parse public key: not an ECDSA public key")
	}

	// Verify that the referenced output's address matches the public key
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	if string(pubKeyBytes) != string(marshaledPubKey) {
		return errors.New("address does not match the public key")
	}

	// Hash the transaction ID
	txIDHash := sha256.Sum256([]byte(tx.ID().String()))

	// Decode the signature
	sigBytes, err := hex.DecodeString(inputTx.Sig)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}
	r := new(big.Int).SetBytes(sigBytes[:len(sigBytes)/2])
	s := new(big.Int).SetBytes(sigBytes[len(sigBytes)/2:])

	// Verify the signature
	if !ecdsa.Verify(publicKey, txIDHash[:], r, s) {
		return fmt.Errorf("invalid signature for transaction input")
	}

	return nil
}

func ValidateTransaction(tx *Transaction, unspent UnspentOutputRepository, blockHeight int) error {
	if err := validateCoinbase(tx, blockHeight); err != nil {
		return err
	}

	for i := 1; i < len(tx.In); i++ {
		if err := validateTransactionIn(tx.In[i], tx, unspent); err != nil {
			return err
		}
	}

	return nil
}
