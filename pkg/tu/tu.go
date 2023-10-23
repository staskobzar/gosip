package tu

type Txn struct {
}

type branch string

type TUPool struct {
	client map[branch]*Txn
	server map[branch]*Txn
}

type TU struct {
	pool TUPool
}

func New() *TU {
	return &TU{
		pool: TUPool{
			client: make(map[branch]*Txn),
			server: make(map[branch]*Txn),
		},
	}
}
