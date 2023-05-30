package forwarders_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Tests the atomicity of cleanup function passed to DeleteForwarder, during DELETE operation
func Test_DeleteForwarder(t *testing.T) {
	t.Parallel()

	_, db := heavyweight.FullTestDBNoFixturesV2(t, "delete_forwarder", nil)
	lggr := logger.TestLogger(t)
	orm := forwarders.NewORM(db, lggr, pgtest.NewQConfig(true))

	addr := testutils.NewAddress()
	chainID := testutils.NewRandomEVMChainID()
	_, err := db.Exec(`INSERT INTO evm_chains (id, created_at, updated_at) VALUES ($1, NOW(), NOW())`, utils.NewBig(chainID))
	require.NoError(t, err)

	fwd, err := orm.CreateForwarder(addr, *utils.NewBig(chainID))
	require.NoError(t, err)
	assert.Equal(t, addr, fwd.Address)

	ErrCleaningUp := errors.New("error during cleanup")

	cleanupCalled := 0

	// Cleanup should fail the first time, causing delete to abort.  When cleanup succeeds the second time,
	//  delete should succeed.  Should fail the 3rd and 4th time since the forwarder has already been deleted.
	//  cleanup should only be called the first two times (when DELETE can succeed).
	rets := []error{ErrCleaningUp, nil, nil, ErrCleaningUp}
	expected := []error{ErrCleaningUp, nil, sql.ErrNoRows, sql.ErrNoRows}

	testCleanupFn := func(q pg.Queryer, evmChainID int64, addr common.Address) error {
		require.Less(t, cleanupCalled, len(rets))
		cleanupCalled++
		return rets[cleanupCalled-1]
	}

	for _, expect := range expected {
		err = orm.DeleteForwarder(fwd.ID, testCleanupFn)
		assert.ErrorIs(t, err, expect)
	}
	assert.Equal(t, 2, cleanupCalled)
}
