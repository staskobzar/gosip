package transaction

import (
	"gosip/pkg/sip"
	"gosip/pkg/transaction/state"
	"gosip/pkg/transaction/timer"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockNonInvite() (*Layer, *sip.Packet) {
	layer := Init()
	pack := &sip.Packet{
		Message: mockRegisterMsg(),
	}
	return layer, pack
}

func TestClientNonInviteInit(t *testing.T) {

	t.Run("set trying and sends first request to transport when terminates no fail", func(t *testing.T) {
		layer, pack := mockNonInvite()
		txn := initClientNonInvite(pack, layer)
		assert.NotNil(t, txn)

		assert.True(t, txn.state.IsTrying())

		assert.Same(t, pack, <-layer.SendTransp())

		assert.NotPanics(t, func() { txn.terminate() })
	})

	t.Run("timer F", func(t *testing.T) {
		setup := func() (*Layer, *sip.Packet) {
			layer, pack := mockNonInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.F = 5 * time.Millisecond
				return t
			}
			return layer, pack
		}

		t.Run("expires and sends timeout error for unreliable transport", func(t *testing.T) {
			layer, pack := setup()
			txn := initClientNonInvite(pack, layer)
			assert.True(t, txn.state.IsTrying())
			assert.Eventually(t, func() bool {
				return txn.state.IsTerminated()
			}, 15*time.Millisecond, 2*time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.ErrorIs(t, ErrTimeout, <-txn.layer.Err())
		})

		t.Run("is not fired for reliable transport", func(t *testing.T) {
			layer, pack := setup()
			pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}

			txn := initClientNonInvite(pack, layer)
			assert.True(t, txn.state.IsTrying())
			assert.Never(t, func() bool {
				return txn.state.IsTerminated()
			}, 8*time.Millisecond, 2*time.Millisecond)
			assert.False(t, txn.state.IsTerminated())
			assert.Zero(t, len(txn.layer.err))
		})
	})

	t.Run("timer E", func(t *testing.T) {
		setup := func() (*Layer, *sip.Packet) {
			layer, pack := mockNonInvite()
			layer.SetupTimers = func(t *timer.Timer) *timer.Timer {
				t.T1 = 2 * time.Millisecond
				t.T2 = 10 * time.Millisecond
				return t
			}
			return layer, pack
		}

		t.Run("not fired for reliable transport", func(t *testing.T) {
			layer, pack := setup()
			pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}

			txn := initClientNonInvite(pack, layer)
			assert.True(t, txn.state.IsTrying())
			assert.Never(t, func() bool {
				return txn.state.IsTerminated()
			}, 18*time.Millisecond, 2*time.Millisecond)

			assert.Equal(t, 1, len(layer.sndTransp))
		})

		t.Run("retransmit on trying state for unreliable transport", func(t *testing.T) {
			layer, pack := setup()

			txn := initClientNonInvite(pack, layer)
			assert.True(t, txn.state.IsTrying())
			assert.Never(t, func() bool {
				return txn.state.IsTerminated()
			}, 18*time.Millisecond, 2*time.Millisecond)

			assert.GreaterOrEqual(t, len(layer.sndTransp), 3)
		})
	})
}

func TestClientNonInviteConsume(t *testing.T) {
	setup := func() (*ClientNonInvite, *sip.Packet, *Layer) {
		layer, pack := mockNonInvite()
		txn := initClientNonInvite(pack, layer)
		return txn, pack, layer
	}
	respPack := func(pack *sip.Packet, code int, reason string) *sip.Packet {
		return &sip.Packet{
			Message: pack.Message.Response(code, reason),
		}
	}

	t.Run("ignore on message is null", func(t *testing.T) {
		txn, _, layer := setup()
		assert.NotPanics(t, func() {
			txn.Consume(&sip.Packet{})
		})
		assert.Len(t, layer.sndTU, 0)
	})

	t.Run("ignore on message is request", func(t *testing.T) {
		txn, pack, layer := setup()
		assert.NotPanics(t, func() {
			txn.Consume(pack)
		})
		assert.Len(t, layer.sndTU, 0)
	})

	t.Run("on trying", func(t *testing.T) {
		t.Run("receive early response enter proceeding", func(t *testing.T) {
			txn, pack, layer := setup()
			resp := respPack(pack, 100, "Trying")
			assert.True(t, txn.state.IsTrying())
			txn.Consume(resp)
			assert.True(t, txn.state.IsProceeding())
			assert.Same(t, resp, <-layer.SendTU())
		})

		t.Run("receive final response enter completed", func(t *testing.T) {
			txn, pack, layer := setup()
			resp := respPack(pack, 200, "Trying")
			assert.True(t, txn.state.IsTrying())
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp, <-layer.SendTU())
		})

		t.Run("receive unknown response", func(t *testing.T) {
			txn, pack, layer := setup()
			resp := respPack(pack, 999, "Foo")
			assert.True(t, txn.state.IsTrying())
			txn.Consume(resp)
			assert.True(t, txn.state.IsTrying())
			assert.Same(t, resp, <-layer.SendTU())
		})
	})

	t.Run("on proceeding", func(t *testing.T) {
		setupProceed := func() (*ClientNonInvite, *sip.Packet, *Layer) {
			txn, pack, layer := setup()
			txn.state.Set(state.Proceeding)
			return txn, pack, layer
		}

		t.Run("receive early response", func(t *testing.T) {
			txn, pack, layer := setupProceed()
			resp := respPack(pack, 100, "Trying")
			assert.True(t, txn.state.IsProceeding())
			txn.Consume(resp)
			assert.True(t, txn.state.IsProceeding())
			assert.Same(t, resp, <-layer.SendTU())
		})

		t.Run("receive final response", func(t *testing.T) {
			txn, pack, layer := setupProceed()
			resp := respPack(pack, 404, "Not Found")
			assert.True(t, txn.state.IsProceeding())
			txn.Consume(resp)
			assert.True(t, txn.state.IsCompleted())
			assert.Same(t, resp, <-layer.SendTU())
		})
	})

	t.Run("on completed", func(t *testing.T) {
		t.Run("absorb responses", func(t *testing.T) {
			txn, pack, layer := setup()
			txn.state.Set(state.Completed)
			assert.True(t, txn.state.IsCompleted())
			txn.Consume(respPack(pack, 100, "Trying"))
			assert.True(t, txn.state.IsCompleted())
			assert.Len(t, layer.sndTU, 0)

			txn.Consume(respPack(pack, 200, "Ok"))
			assert.True(t, txn.state.IsCompleted())
			assert.Len(t, layer.sndTU, 0)
		})

		t.Run("timer K fired for unreliable transport", func(t *testing.T) {
			txn, _, layer := setup()
			txn.timer.T4 = 10 * time.Millisecond
			txn.completed()
			assert.True(t, txn.state.IsCompleted())
			assert.Eventually(t, func() bool {
				return txn.state.IsTerminated()
			}, 15*time.Millisecond, 2*time.Millisecond)
			assert.True(t, txn.state.IsTerminated())
			assert.Len(t, layer.sndTU, 0)
		})

		t.Run("timer K is not used for reliable transport", func(t *testing.T) {
			txn, pack, layer := setup()
			pack.SendTo = []net.Addr{&net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 5060}}
			txn.completed()
			assert.True(t, txn.state.IsTerminated())
			assert.Len(t, layer.sndTU, 0)
		})
	})
}
