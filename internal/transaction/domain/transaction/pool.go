package transaction

type PoolRepository interface {
	Add(*Transaction) error
	Exists(ID) bool
	Remove(...ID) error
	GetAll() []*Transaction
}

type Pool struct {
	pool PoolRepository
}

func NewPool(pool PoolRepository) *Pool {
	return &Pool{
		pool: pool,
	}
}

func (p *Pool) Add(tx *Transaction) error {
	// TODO#25 validation

	return p.pool.Add(tx)
}

func (p *Pool) Exists(id ID) bool {
	return p.pool.Exists(id)
}

func (p *Pool) Update(unspent []UnspentOutput) error {
	unspentIDs := make(map[ID]struct{})
	for _, u := range unspent {
		unspentIDs[u.outputID] = struct{}{}
	}
	spentIDs := make(map[ID]struct{})
	current := p.pool.GetAll()

	for _, u := range current {
		for _, i := range u.Inputs() {
			if _, ok := unspentIDs[i.outputID]; !ok {
				spentIDs[i.outputID] = struct{}{}
				break
			}
		}
	}

	toRemove := make([]ID, len(spentIDs))
	i := 0
	for id := range spentIDs {
		toRemove[i] = id
		i++
	}

	return p.pool.Remove(toRemove...)
}
