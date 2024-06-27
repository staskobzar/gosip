// Package pool provides storage for transactions
package pool

import (
	"errors"
	"fmt"
	"gosip/pkg/logger"
	"gosip/pkg/sip"
	"gosip/pkg/sipmsg"
	"sync"
)

// Error for pool package
var Error = errors.New("transactions:pool")

// Pool storage of the SIP transactions
type Pool struct {
	mu sync.RWMutex
	m  map[string]sip.Transaction
}

// New pool create
func New() *Pool {
	return &Pool{
		m: make(map[string]sip.Transaction),
	}
}

// Add a transaction to the pool
func (p *Pool) Add(txn sip.Transaction) error {
	if _, ok := p.Get(txn.BranchID()); ok {
		return fmt.Errorf("%w: already exists %q", Error, txn.BranchID())
	}
	p.push(txn)
	logger.Log("txn:pool: number of transactions in the pool after add: %d", p.Len())
	return nil
}

// Delete a transaction from the pool by id
func (p *Pool) Delete(id string) {
	p.del(id)
	logger.Log("txn:pool: number of transactions in the pool after delete: %d", p.Len())
}

// Get transcation from pool by top Via branch and returns
// the transaction if found and true if found or nil and false if not found
func (p *Pool) Get(branch string) (sip.Transaction, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	txn, ok := p.m[branch]
	return txn, ok
}

// Len returns a number of transactions in the pool
func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.m)
}

// Match transaction in pools against SIP message
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
