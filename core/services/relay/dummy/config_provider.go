package dummy

import (
	"context"
	"encoding/json"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type RelayConfig struct {
	// TODO: need to map json relay cfg into config provider somehow
}

type configProvider struct {
	lggr logger.Logger
	cfg  RelayConfig

	digester ocrtypes.OffchainConfigDigester
	tracker  ocrtypes.ContractConfigTracker
}

func NewConfigProvider(lggr logger.Logger, rargs types.RelayArgs) (types.ConfigProvider, error) {
	cp := &configProvider{lggr: lggr.Named("DummyConfigProvider")}
	err := json.Unmarshal(rargs.RelayConfig, &cp.cfg)
	if err != nil {
		return nil, err
	}
	cp.digester, err = NewOffchainConfigDigester(cp.lggr, cp.cfg)
	if err != nil {
		return nil, err
	}
	cp.tracker = NewContractConfigTracker(cp.lggr, cp.cfg)
	return cp, nil
}

func (cp *configProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return cp.digester
}
func (cp *configProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker { return cp.tracker }
func (cp *configProvider) Name() string                                          { return cp.lggr.Name() }
func (*configProvider) Start(context.Context) error                              { return nil }
func (*configProvider) Close() error                                             { return nil }
func (*configProvider) Ready() error                                             { return nil }
func (cp *configProvider) HealthReport() map[string]error {
	return map[string]error{cp.lggr.Name(): nil}
}
