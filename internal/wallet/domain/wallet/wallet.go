package wallet

type ID string

type Algorithm int

const (
	RSA Algorithm = iota
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
