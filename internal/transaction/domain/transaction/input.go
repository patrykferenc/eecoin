package transaction

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

type Input struct {
	outputID    ID
	outputIndex int
	signature   string // TODO#30 signature struct
}

func (i *Input) sign(signer crypto.Signer, idToSign ID, referencedOutput UnspentOutput) error {
	ourAddress, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		fmt.Errorf("error marshalling public key: %w", err)
	}
	if string(ourAddress) != referencedOutput.address {
		return fmt.Errorf("output address does not match the signer address")
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
