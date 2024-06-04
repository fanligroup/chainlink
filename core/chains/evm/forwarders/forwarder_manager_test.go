package forwarders_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

var GetAuthorisedSendersABI = evmtypes.MustGetABI(authorized_receiver.AuthorizedReceiverABI).Methods["getAuthorizedSenders"]

var SimpleOracleCallABI = evmtypes.MustGetABI(operator_wrapper.OperatorABI).Methods["getChainlinkToken"]

func TestFwdMgr_MaybeForwardTransaction(t *testing.T) {
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	owner := testutils.MustNewSimTransactor(t)
	ctx := testutils.Context(t)

	b := simulated.NewBackend(types.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	t.Cleanup(func() { b.Close() })
	linkAddr := common.HexToAddress("0x01BE23585060835E02B77ef475b0Cc51aA1e0709")
	operatorAddr, _, _, err := operator_wrapper.DeployOperator(owner, b.Client(), linkAddr, owner.From)
	require.NoError(t, err)
	forwarderAddr, _, forwarder, err := authorized_forwarder.DeployAuthorizedForwarder(owner, b.Client(), linkAddr, owner.From, operatorAddr, []byte{})
	require.NoError(t, err)
	b.Commit()
	_, err = forwarder.SetAuthorizedSenders(owner, []common.Address{owner.From})
	require.NoError(t, err)
	b.Commit()
	authorized, err := forwarder.GetAuthorizedSenders(nil)
	require.NoError(t, err)
	t.Log(authorized)

	evmClient := client.NewSimulatedBackendClient(t, b, testutils.FixtureChainID)

	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), evmClient, lggr, lpOpts)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	fwdMgr.ORM = forwarders.NewORM(db)

	fwd, err := fwdMgr.ORM.CreateForwarder(ctx, forwarderAddr, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(ctx, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	require.NoError(t, fwdMgr.Start(testutils.Context(t)))
	addr, err := fwdMgr.ForwarderFor(owner.From)
	require.NoError(t, err)
	require.Equal(t, addr.String(), forwarderAddr.String())
	err = fwdMgr.Close()
	require.NoError(t, err)

	cleanupCalled := false
	cleanup := func(tx sqlutil.DataSource, evmChainId int64, addr common.Address) error {
		require.Equal(t, testutils.FixtureChainID.Int64(), evmChainId)
		require.Equal(t, forwarderAddr, addr)
		require.NotNil(t, tx)
		cleanupCalled = true
		return nil
	}

	err = fwdMgr.ORM.DeleteForwarder(ctx, fwd.ID, cleanup)
	assert.NoError(t, err)
	assert.True(t, cleanupCalled)
}

func TestFwdMgr_AccountUnauthorizedToForward_SkipsForwarding(t *testing.T) {
	lggr := logger.Test(t)
	db := pgtest.NewSqlxDB(t)
	ctx := testutils.Context(t)
	cfg := configtest.NewTestGeneralConfig(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	owner := testutils.MustNewSimTransactor(t)
	b := simulated.NewBackend(types.GenesisAlloc{
		owner.From: {
			Balance: big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
	}, simulated.WithBlockGasLimit(10e6))
	t.Cleanup(func() { b.Close() })
	linkAddr := common.HexToAddress("0x01BE23585060835E02B77ef475b0Cc51aA1e0709")
	operatorAddr, _, _, err := operator_wrapper.DeployOperator(owner, b.Client(), linkAddr, owner.From)
	require.NoError(t, err)

	forwarderAddr, _, _, err := authorized_forwarder.DeployAuthorizedForwarder(owner, b.Client(), linkAddr, owner.From, operatorAddr, []byte{})
	require.NoError(t, err)
	b.Commit()

	evmClient := client.NewSimulatedBackendClient(t, b, testutils.FixtureChainID)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), evmClient, lggr, lpOpts)
	fwdMgr := forwarders.NewFwdMgr(db, evmClient, lp, lggr, evmcfg.EVM())
	fwdMgr.ORM = forwarders.NewORM(db)

	_, err = fwdMgr.ORM.CreateForwarder(ctx, forwarderAddr, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	lst, err := fwdMgr.ORM.FindForwardersByChain(ctx, ubig.Big(*testutils.FixtureChainID))
	require.NoError(t, err)
	require.Equal(t, len(lst), 1)
	require.Equal(t, lst[0].Address, forwarderAddr)

	err = fwdMgr.Start(testutils.Context(t))
	require.NoError(t, err)
	addr, err := fwdMgr.ForwarderFor(owner.From)
	require.ErrorContains(t, err, "Cannot find forwarder for given EOA")
	require.True(t, utils.IsZero(addr))
	err = fwdMgr.Close()
	require.NoError(t, err)
}
