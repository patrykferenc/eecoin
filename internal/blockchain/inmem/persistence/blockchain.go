package persistence

import (
	"encoding/json"
	"os"

	bc "github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type ChainDto struct {
	Blocks []blockDTO `json:"blocks"`
}

func MapToDto(blockchain bc.BlockChain) ChainDto {
	dtoBlocks := make([]blockDTO, len(blockchain.Blocks))
	for i, block := range blockchain.Blocks {
		dtoBlocks[i] = asDTO(block)
	}
	return ChainDto{Blocks: dtoBlocks}
}

func MapToActual(chain ChainDto) (bc.BlockChain, error) {
	dtoBlocks := make([]bc.Block, len(chain.Blocks))
	for i, block := range chain.Blocks {

		transactions := make([]transaction.Transaction, len(block.Transactions))
		for i, trscnion := range block.Transactions {
			translated, err := asModel(trscnion)
			if err != nil {
				return bc.BlockChain{}, err
			}
			transactions[i] = *translated
		}

		dtoBlocks[i] = bc.Block{
			Index:          block.Index,
			TimestampMilis: block.TimestampMilis,
			ContentHash:    block.ContentHash,
			PrevHash:       block.PrevHash,
			Transactions:   transactions,
			Challenge:      challengeDTOToModel(block.Challange),
		}
	}
	output, err := bc.ImportBlockchain(dtoBlocks)
	if err != nil {
		return bc.BlockChain{}, err
	}
	return *output, nil
}

func Persist(chain bc.BlockChain, path string) error {
	mappedToDto := MapToDto(chain)
	b, err := json.Marshal(mappedToDto)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

func Load(path string) (*bc.BlockChain, error) {
	persistedContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	blockchainDto := ChainDto{}
	err = json.Unmarshal(persistedContent, &blockchainDto)
	if err != nil {
		return nil, err
	}

	blockchain, err := MapToActual(blockchainDto)
	if err != nil {
		return nil, err
	}
	return &blockchain, nil
}

type blockDTO struct {
	Index          int              `json:"index"`
	TimestampMilis int64            `json:"timestamp"`
	ContentHash    string           `json:"content_hash"`
	PrevHash       string           `json:"prev_hash"`
	Transactions   []transactionDTO `json:"transactions"` // TODO#30
	Challange      challengeDTO     `json:"challenge"`
}

func asDTO(block bc.Block) blockDTO {
	transactions := make([]transactionDTO, len(block.Transactions))
	for i, trscnion := range block.Transactions {
		transactions[i] = transDTO(trscnion)
	}

	return blockDTO{
		Index:          block.Index,
		TimestampMilis: block.TimestampMilis,
		ContentHash:    block.ContentHash,
		PrevHash:       block.PrevHash,
		Transactions:   transactions,
		Challange:      challengeModelToDTO(block.Challenge),
	}
}

type inputDTO struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Signature   string `json:"signature"`
}

func (i inputDTO) asInput() *transaction.Input {
	return transaction.NewInput(transaction.ID(i.OutputID), i.OutputIndex, i.Signature)
}

type outputDTO struct {
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

func (o outputDTO) asOutput() *transaction.Output {
	return transaction.NewOutput(o.Amount, o.Address)
}

type transactionDTO struct {
	ID      string      `json:"id"`
	Inputs  []inputDTO  `json:"inputs"`
	Outputs []outputDTO `json:"outputs"`
}

func transDTO(tx transaction.Transaction) transactionDTO {
	inputs := make([]inputDTO, len(tx.Inputs()))
	for i, in := range tx.Inputs() {
		inputs[i] = inputDTO{
			OutputID:    in.OutputID().String(),
			OutputIndex: in.OutputIndex(),
			Signature:   in.Signature(),
		}
	}

	outputs := make([]outputDTO, len(tx.Outputs()))
	for i, out := range tx.Outputs() {
		outputs[i] = outputDTO{
			Amount:  out.Amount(),
			Address: out.Address(),
		}
	}

	return transactionDTO{
		ID:      tx.ID().String(),
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func asModel(dto transactionDTO) (*transaction.Transaction, error) {
	inputs := make([]*transaction.Input, len(dto.Inputs))
	for i, in := range dto.Inputs {
		inputs[i] = in.asInput()
	}

	outputs := make([]*transaction.Output, len(dto.Outputs))
	for i, out := range dto.Outputs {
		outputs[i] = out.asOutput()
	}

	return transaction.NewFrom(inputs, outputs)
}

type challengeDTO struct {
	Difficulty    int    `json:"difficulty"`
	Nonce         uint32 `json:"nonce"`
	HashValue     string `json:"hash_value"`
	TimeCapMillis int64  `json:"time_cap_millis"`
}

func challengeDTOToModel(dto challengeDTO) bc.Challenge {
	return bc.Challenge{
		Difficulty:    dto.Difficulty,
		Nonce:         dto.Nonce,
		HashValue:     dto.HashValue,
		TimeCapMillis: dto.TimeCapMillis,
	}
}

func challengeModelToDTO(challenge bc.Challenge) challengeDTO {
	return challengeDTO{
		Difficulty:    challenge.Difficulty,
		Nonce:         challenge.Nonce,
		HashValue:     challenge.HashValue,
		TimeCapMillis: challenge.TimeCapMillis,
	}
}
