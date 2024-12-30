package transaction

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

type Input struct {
	OutputId  ID
	OutputIdx int
	Sig       string
}

func NewInput(outputID ID, outputIndex int, signature string) *Input {
	return &Input{
		OutputId:  outputID,
		OutputIdx: outputIndex,
		Sig:       signature,
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

	i.Sig = string(s)
	return nil
}

func (i Input) Signature() string {
	return i.Sig
}

func (i Input) OutputID() ID {
	return i.OutputId
}

func (i Input) OutputIndex() int {
	return i.OutputIdx
}
