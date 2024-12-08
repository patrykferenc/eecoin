package blockchain

const (
	BlockGenerationDefaultIntervalMillis int64 = 100 * 60
	DifficultyAdjustmentInterval               = 10
)

func GetDifficulty(chain BlockChain) (int, error) {
	latest := chain.GetLast()
	if latest.Index != 0 && latest.Index%DifficultyAdjustmentInterval == 0 {
		calculatedDifficulty, err := getAdjustedDifficulty(chain)
		return calculatedDifficulty, err
	}
	return latest.Challenge.Difficulty, nil
}

func getAdjustedDifficulty(chain BlockChain) (int, error) {
	previouslyAdjusted, err := chain.GetBlock(len(chain.Blocks) - DifficultyAdjustmentInterval)
	if err != nil {
		return 0, err
	}

	expectedTime := BlockGenerationDefaultIntervalMillis * DifficultyAdjustmentInterval
	actualTime := chain.GetLast().TimestampMilis - previouslyAdjusted.TimestampMilis

	lower := ifThen(previouslyAdjusted.Challenge.Difficulty-1 <= 0, previouslyAdjusted.Challenge.Difficulty, previouslyAdjusted.Challenge.Difficulty-1)
	higher := ifThen(previouslyAdjusted.Challenge.Difficulty+1 >= 256, previouslyAdjusted.Challenge.Difficulty, previouslyAdjusted.Challenge.Difficulty+1)

	switch {
	case actualTime < expectedTime/2:
		return higher, nil
	case actualTime > expectedTime*2:
		return lower, nil
	}
	return chain.GetLast().Challenge.Difficulty, nil
}

func ifThen[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
