package txmgr

import (
	"context"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/common/txmgr/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var _ txmgrtypes.TxStrategy = SendEveryStrategy{}

// NewQueueingTxStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded if a queue size is specified, and otherwise does not drop transactions.
func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) (strategy txmgrtypes.TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize, queryTimeout)
	} else {
		strategy = SendEveryStrategy{}
	}
	return
}

// NewSendEveryStrategy creates a new TxStrategy that does not drop transactions.
func NewSendEveryStrategy() txmgrtypes.TxStrategy {
	return SendEveryStrategy{}
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct{}

func (SendEveryStrategy) Subject() uuid.NullUUID                              { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(pruneService any, opt any) (int64, error) { return 0, nil }

var _ types.TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy struct {
	subject      uuid.UUID
	queueSize    uint32
	queryTimeout time.Duration
}

// NewDropOldestStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded.
func NewDropOldestStrategy(subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) DropOldestStrategy {
	return DropOldestStrategy{subject, queueSize, queryTimeout}
}

func (s DropOldestStrategy) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy) PruneQueue(pruneService any, opt any) (n int64, err error) {
	qopt, err := ToQOpt(opt)
	if err != nil {
		return 0, errors.Wrap(err, "DropOldestStrategy#PruneQueue failed")
	}

	orm, ok := pruneService.(ORM)
	if !ok {
		return 0, errors.Wrap(err, "DropOldestStrategy#PruneQueue failed invalid pruneService")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)

	defer cancel()
	n, err = orm.PruneUnstartedEthTxQueue(s.queueSize, s.subject, pg.WithParentCtx(ctx), qopt)
	if err != nil {
		return 0, errors.Wrap(err, "DropOldestStrategy#PruneQueue failed")
	}
	return
}
