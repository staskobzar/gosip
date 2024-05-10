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

type Pool struct {
	mu sync.RWMutex
	m  map[string]sip.Transaction
}

func New() *Pool {
	return &Pool{
		m: make(map[string]sip.Transaction),
	}
}

func (p *Pool) Add(txn sip.Transaction) error {
	if _, ok := p.Get(txn.BranchID()); ok {
		return fmt.Errorf("%w: already exists %q", Error, txn.BranchID())
	}
	p.push(txn)
	logger.Log("txn:pool: number of transactions in the pool after add: %d", p.Len())
	return nil
}

func (p *Pool) Delete(id string) {
	p.del(id)
	logger.Log("txn:pool: number of transactions in the pool after delete: %d", p.Len())
}

func (p *Pool) Get(branch string) (sip.Transaction, bool) {
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

func (p *Pool) Match(msg *sipmsg.Message) (sip.Transaction, bool) {
	branch := msg.TopViaBranch()
	txn, found := p.Get(branch)
	if !found {
		logger.Wrn("txn:pool: no transaction found for branch %q", branch)
		return nil, false
	}

	return txn.Match(msg)
}

func (p *Pool) push(txn sip.Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.m[txn.BranchID()] = txn
}

func (p *Pool) del(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.m, id)
}
