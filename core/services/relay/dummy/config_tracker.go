package dummy

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type configTracker struct{}

func NewContractConfigTracker(lggr logger.Logger, cfg RelayConfig) ocrtypes.ContractConfigTracker {
	return &configTracker{}
}

// Notify may optionally emit notification events when the contract's
// configuration changes. This is purely used as an optimization reducing
// the delay between a configuration change and its enactment. Implementors
// who don't care about this may simply return a nil channel.
//
// The returned channel should never be closed.
func (ct *configTracker) Notify() <-chan struct{} {
	return nil
}

// LatestConfigDetails returns information about the latest configuration,
// but not the configuration itself.
func (ct *configTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	return 0, ocrtypes.ConfigDigest{}, nil
}

// LatestConfig returns the latest configuration.
func (ct *configTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ocrtypes.ContractConfig{}, nil
}

// LatestBlockHeight returns the height of the most recent block in the chain.
func (ct *configTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	return 0, nil
}
