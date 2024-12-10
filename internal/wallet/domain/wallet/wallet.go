package wallet

import "errors"

type ID string

type Algorithm int

const (
	RSA Algorithm = iota
	ECDSA
)

var (
	ErrPrivateKeyNotFound = errors.New("private Key not found")
	NoKeysFound           = errors.New("no keys found")
	PemParseError         = errors.New("pem parse error")
)

type Wallet[T any, E any] interface {
	SetMainIdentity(Key[T, E]) error
	Add(Key[T, E]) error
	Type() Algorithm
}

type Key[T any, E any] struct {
	private T
	Public  E
	algType Algorithm
}

func (k *Key[T, E]) Private() T {
	return k.private
}

type KeyElement[T any] struct {
	Key     *T
	Present bool
}
