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

func (p *Pool) Update(unspent []UnspentOutput) {
	// TODO#25 update pool
}
