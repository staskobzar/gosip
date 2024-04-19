package pool

import (
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"sync"
)

var Error = errors.New("transactions:pool")

type Transaction interface {
	BranchID() string
	Consume(*sip.Packet)
	Match(msg *sipmsg.Message) (Transaction, bool)
}

type Pool struct {
	mu sync.RWMutex
	m  map[string]Transaction
}

func New() *Pool {
	return &Pool{
		m: make(map[string]Transaction),
	}
}

func (p *Pool) Add(txn Transaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.m[txn.BranchID()]; ok {
		return fmt.Errorf("%w: already exists %q", Error, txn.BranchID())
	}
	p.m[txn.BranchID()] = txn
	// do not use p.Len to avoid deadlock
	logger.Log("txn:pool: number of transactions in the pool after add: %d", len(p.m))
	return nil
}

func (p *Pool) Delete(txn Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.m, txn.BranchID())
	logger.Log("txn:pool: number of transactions in the pool after delete: %d", len(p.m))
}

func (p *Pool) Get(branch string) (Transaction, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	txn, ok := p.m[branch]
	return txn, ok
}

func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.m)
}

func (p *Pool) Match(msg *sipmsg.Message) (Transaction, bool) {
	branch := msg.TopViaBranch()
	txn, found := p.Get(branch)
	if !found {
		logger.Wrn("txn:pool: no transaction found for branch %q", branch)
		return nil, false
	}

	return txn.Match(msg)
}