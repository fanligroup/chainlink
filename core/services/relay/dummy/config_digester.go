package dummy

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type configDigester struct{}

func NewOffchainConfigDigester(lggr logger.Logger, cfg RelayConfig) (ocrtypes.OffchainConfigDigester, error) {
	return &configDigester{}, nil
}

// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
func (cd *configDigester) ConfigDigest(ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	return ocrtypes.ConfigDigest{}, nil
}

// This should return the same constant value on every invocation
func (cd *configDigester) ConfigDigestPrefix() (ocrtypes.ConfigDigestPrefix, error) {
	return 0, nil
}
