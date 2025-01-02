package transaction

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

type Input struct {
	outputID    ID
	outputIndex int
	signature   string
}

func NewInput(outputID ID, outputIndex int, signature string) *Input {
	return &Input{
		outputID:    outputID,
		outputIndex: outputIndex,
		signature:   signature,
	}
}

func (i *Input) sign(signer crypto.Signer, idToSign ID, referencedOutput UnspentOutput) error {
	ourAddress, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return fmt.Errorf("error marshalling public key: %w", err)
	}
	if hex.EncodeToString(ourAddress) != referencedOutput.address {
		return fmt.Errorf("output Addr does not match the signer Addr")
	}

	s, err := signer.Sign(rand.Reader, []byte(idToSign), crypto.SHA256)
	if err != nil {
		return fmt.Errorf("error signing input: %w", err)
	}

	i.signature = string(s)
	return nil
}

func (i Input) Signature() string {
	return i.signature
}

func (i Input) OutputID() ID {
	return i.outputID
}

func (i Input) OutputIndex() int {
	return i.outputIndex
}

func (i Input) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%d%s", i.outputID, i.outputIndex, i.signature)), nil
}
