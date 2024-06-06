package dummy

import (
	"context"
	"encoding/json"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type channelDefinitionCache struct {
	lggr logger.Logger

	definitions llotypes.ChannelDefinitions
}

func NewChannelDefinitionCache(lggr logger.Logger, dfns string) (llotypes.ChannelDefinitionCache, error) {
	var definitions llotypes.ChannelDefinitions
	err := json.Unmarshal([]byte(dfns), &definitions)
	if err != nil {
		return nil, err
	}
	return &channelDefinitionCache{lggr, definitions}, nil
}

func (cdc *channelDefinitionCache) Definitions() llotypes.ChannelDefinitions { return cdc.definitions }
func (cdc *channelDefinitionCache) Start(context.Context) error              { return nil }
func (cdc *channelDefinitionCache) Close() error                             { return nil }
func (cdc *channelDefinitionCache) Ready() error                             { return nil }
func (cdc *channelDefinitionCache) HealthReport() map[string]error           { return nil }
func (cdc *channelDefinitionCache) Name() string                             { return cdc.lggr.Name() }
