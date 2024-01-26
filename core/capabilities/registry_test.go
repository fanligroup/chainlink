package capabilities

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type cmockCapability struct {
	CapabilityInfo
}

func (m *cmockCapability) Execute(ctx context.Context, callback chan CapabilityResponse, inputs values.Map) error {
	return nil
}

func TestRegistry(t *testing.T) {
	ctx := context.Background()

	r := NewRegistry()

	id := "capability-1"
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &cmockCapability{CapabilityInfo: ci}
	err = r.Add(ctx, c)
	require.NoError(t, err)

	gc, err := r.Get(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, c, gc)

	cs := r.List(ctx)
	assert.Len(t, cs, 1)
	assert.Equal(t, c, cs[0])
}

func TestRegistry_NoDuplicateIDs(t *testing.T) {
	ctx := context.Background()
	r := NewRegistry()

	id := "capability-1"
	ci, err := NewCapabilityInfo(
		id,
		CapabilityTypeAction,
		"capability-1-description",
		"v1.0.0",
	)
	require.NoError(t, err)

	c := &cmockCapability{CapabilityInfo: ci}
	err = r.Add(ctx, c)
	require.NoError(t, err)

	ci, err = NewCapabilityInfo(
		id,
		CapabilityTypeConsensus,
		"capability-2-description",
		"v1.0.0",
	)
	require.NoError(t, err)
	c2 := &cmockCapability{CapabilityInfo: ci}

	err = r.Add(ctx, c2)
	assert.ErrorContains(t, err, "capability with id: capability-1 already exists")
}

func TestRegistry_ChecksExecutionAPIByType(t *testing.T) {
	tcs := []struct {
		name          string
		newCapability func() BaseCapability
		errContains   string
	}{
		{
			name: "action, sync",
			newCapability: func() BaseCapability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeAction,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "target, sync",
			newCapability: func() BaseCapability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTarget,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "trigger, async",
			newCapability: func() BaseCapability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeTrigger,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
		{
			name: "reports, async",
			newCapability: func() BaseCapability {
				id := uuid.New().String()
				ci, err := NewCapabilityInfo(
					id,
					CapabilityTypeConsensus,
					"capability-1-description",
					"v1.0.0",
				)
				require.NoError(t, err)

				return &cmockCapability{CapabilityInfo: ci}
			},
		},
	}

	ctx := context.Background()
	reg := NewRegistry()
	for _, tc := range tcs {
		c := tc.newCapability()
		err := reg.Add(ctx, c)
		require.NoError(t, err)
	}
}
