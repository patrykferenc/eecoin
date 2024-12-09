package transaction

type Input struct {
	outputID    ID
	outputIndex int
	signature   string // TODO#30 signature struct
}
