package blockchain

import "errors"

var NotFoundAmongLocalChainCopies = errors.New("block not found among local chain copies")

type Mediator interface {
	AddBlock(Block) error
	GetPrimaryChain() (*BlockChain, error)
}

type ForkMediator struct {
	primaryChain *BlockChain
	chains       []*BlockChain
}

func NewForkMediator(primaryChain *BlockChain) *ForkMediator {
	chains := make([]*BlockChain, 0)
	chains = append(chains, primaryChain)
	return &ForkMediator{
		primaryChain: primaryChain,
		chains:       chains,
	}
}

func (f *ForkMediator) AddBlock(new Block) error {
	primaryChain := f.primaryChain
	err := primaryChain.AddBlock(new)
	if err == nil {
		return nil
	}
	// assuming we've already asked other nodes for their local copies which are not the primary chain
	if errors.Is(err, PossibleForkBlockchain) {
		var addedToAnyChain bool
		var maxCumulativeDifficulty = f.primaryChain.GetCumulativeDifficulty()
		var maxChainPtr = f.primaryChain
		// reverse iteration - the most recent will be probably added last, but we assume the earliest one is the truest one
		// so the first chain in the slice should have the highest priority to be the primary chain
		for i := len(f.chains) - 1; i >= 0; i-- {
			err := f.chains[i].AddBlock(new)
			if err == nil && f.chains[i].GetCumulativeDifficulty() > maxCumulativeDifficulty {
				maxCumulativeDifficulty = f.chains[i].GetCumulativeDifficulty()
				maxChainPtr = f.chains[i]
				addedToAnyChain = true
			}
			if errors.Is(err, PossibleForkBlockchain) {
				// copy the chain up to new block index and add it to chains
				detectedForkChain := f.chains[i].GetShortenedUpToIndex(new.Index)
				err := detectedForkChain.AddBlock(new)
				if err == nil {
					if detectedForkChain.GetCumulativeDifficulty() > maxCumulativeDifficulty {
						f.primaryChain = &detectedForkChain
						f.chains[0] = &detectedForkChain
						f.chains = f.chains[:1]
						addedToAnyChain = true
					} else if detectedForkChain.GetCumulativeDifficulty() == maxCumulativeDifficulty {
						f.chains = append(f.chains, &detectedForkChain)
						addedToAnyChain = true
					}
				}
			}
		}

		if !addedToAnyChain {
			return NotFoundAmongLocalChainCopies
		}

		f.cleanupChainsWithLowerThanPrimaryDifficulty(maxCumulativeDifficulty)
		primaryChain = maxChainPtr
	}
	return nil
}

func (f *ForkMediator) cleanupChainsWithLowerThanPrimaryDifficulty(maxCumulativeDifficulty int64) {
	for i, chain := range f.chains {
		if chain.GetCumulativeDifficulty() < maxCumulativeDifficulty {
			f.chains = append(f.chains[:i], f.chains[i+1:]...)
		}
	}
}

func (f *ForkMediator) GetPrimaryChain() (*BlockChain, error) {
	if f.primaryChain == nil {
		return nil, errors.New("no primary chain")
	}
	return f.primaryChain, nil
}
