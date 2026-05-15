package app

import (
	"context"
	"encoding/json"
	"time"

	"cosmossdk.io/math"
	"fmt"
	"io"
	"os"
	"path/filepath"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// IBC (v10 — no capability module)
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcporttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	// Custom modules
	"syreen/x/dex"
	dexkeeper "syreen/x/dex/keeper"
	dextypes "syreen/x/dex/types"
	"syreen/x/abstractaccount"
	aaante "syreen/x/abstractaccount/ante"
	abstractaccountkeeper "syreen/x/abstractaccount/keeper"
	abstractaccounttypes "syreen/x/abstractaccount/types"
	"syreen/x/feemarket"
	feemarketkeeper "syreen/x/feemarket/keeper"
	feemarkettypes "syreen/x/feemarket/types"
	"syreen/x/mevprotection"
	mevprotectionkeeper "syreen/x/mevprotection/keeper"
	mevprotectiontypes "syreen/x/mevprotection/types"
	"syreen/x/compute"
	computekeeper "syreen/x/compute/keeper"
	computetypes "syreen/x/compute/types"
	"syreen/x/intent"
	intentkeeper "syreen/x/intent/keeper"
	intenttypes "syreen/x/intent/types"
	"syreen/x/tokenfactory"
	tokenfactorykeeper "syreen/x/tokenfactory/keeper"
	tokenfactorytypes "syreen/x/tokenfactory/types"
	"syreen/x/evm"
	evmkeeper "syreen/x/evm/keeper"
	evmtypes "syreen/x/evm/types"
	"syreen/x/identity"
	identitykeeper "syreen/x/identity/keeper"
	identitytypes "syreen/x/identity/types"
	"syreen/x/payments"
	paymentskeeper "syreen/x/payments/keeper"
	paymentstypes "syreen/x/payments/types"
	"syreen/x/farming"
	farmingkeeper "syreen/x/farming/keeper"
	farmingtypes "syreen/x/farming/types"
	"syreen/x/perps"
	perpskeeper "syreen/x/perps/keeper"
	perpstypes "syreen/x/perps/types"
	"syreen/x/lending"
	lendingkeeper "syreen/x/lending/keeper"
	lendingtypes "syreen/x/lending/types"
	"syreen/x/predict"
	predictkeeper "syreen/x/predict/keeper"
	predicttypes "syreen/x/predict/types"
	"syreen/x/launchpad"
	launchpadkeeper "syreen/x/launchpad/keeper"
	launchpadtypes "syreen/x/launchpad/types"
	"syreen/x/clmm"
	clmmkeeper "syreen/x/clmm/keeper"
	clmmtypes "syreen/x/clmm/types"
	"syreen/x/vault"
	vaultkeeper "syreen/x/vault/keeper"
	vaulttypes "syreen/x/vault/types"
	"syreen/x/flashloan"
	flashloankeeper "syreen/x/flashloan/keeper"
	flashloantypes "syreen/x/flashloan/types"
	"syreen/x/portfolio"
	portfoliokeeper "syreen/x/portfolio/keeper"
	portfoliotypes "syreen/x/portfolio/types"
	"syreen/x/options"
	optionskeeper "syreen/x/options/keeper"
	optionstypes "syreen/x/options/types"
	"syreen/x/aiagent"
	aiagentkeeper "syreen/x/aiagent/keeper"
	aiagenttypes "syreen/x/aiagent/types"

	// Parallel execution engine
	"syreen/engine/parallel"

	// JSON-RPC server
	"syreen/server/jsonrpc"

	"github.com/spf13/cast"
)

const (
	AppName          = "SyreenChain"
	Bech32MainPrefix = "syreen"
	CoinType         = 118
)

var (
	DefaultNodeHome string

	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		tokenfactorytypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		feemarkettypes.ModuleName:        {authtypes.Burner},
		abstractaccounttypes.ModuleName:  nil,
		intenttypes.ModuleName:           nil,
		computetypes.ModuleName:          nil,
		evmtypes.ModuleName:             {authtypes.Minter, authtypes.Burner},
		dextypes.ModuleName:              {authtypes.Minter, authtypes.Burner},
		identitytypes.ModuleName:         nil,
		paymentstypes.ModuleName:         nil,
		farmingtypes.ModuleName:          nil,
		perpstypes.ModuleName:            nil,
		lendingtypes.ModuleName:          nil,
		predicttypes.ModuleName:          nil,
		launchpadtypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		clmmtypes.ModuleName:             nil,
		flashloantypes.ModuleName:        nil,
		vaulttypes.ModuleName:            nil,
		portfoliotypes.ModuleName:        nil,
		optionstypes.ModuleName:          nil,
		aiagenttypes.ModuleName:          nil,
		mevprotectiontypes.ModuleName:    {authtypes.Minter},
	}
)

var (
	_ runtime.AppI            = (*SyreenApp)(nil)
	_ servertypes.Application = (*SyreenApp)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".syreen")

	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(Bech32MainPrefix, Bech32MainPrefix+"pub")
	sdkConfig.SetBech32PrefixForValidator(Bech32MainPrefix+"valoper", Bech32MainPrefix+"valoperpub")
	sdkConfig.SetBech32PrefixForConsensusNode(Bech32MainPrefix+"valcons", Bech32MainPrefix+"valconspub")
	sdkConfig.SetCoinType(CoinType)
	sdkConfig.Seal()
}

type SyreenApp struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// Core keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// Additional SDK keepers
	EvidenceKeeper  evidencekeeper.Keeper
	FeegrantKeeper  feegrantkeeper.Keeper
	AuthzKeeper     authzkeeper.Keeper

	// IBC keepers (v10 — no capability, no scoped keepers)
	IBCKeeper      *ibckeeper.Keeper
	TransferKeeper ibctransferkeeper.Keeper

	// Custom keepers
	TokenfactoryKeeper    tokenfactorykeeper.Keeper
	FeeMarketKeeper       feemarketkeeper.Keeper
	MEVProtectionKeeper   mevprotectionkeeper.Keeper
	AbstractAccountKeeper abstractaccountkeeper.Keeper
	IntentKeeper          *intentkeeper.Keeper
	ComputeKeeper         computekeeper.Keeper
	EVMKeeper             evmkeeper.Keeper
	DexKeeper             *dexkeeper.Keeper
	IdentityKeeper        *identitykeeper.Keeper
	PaymentsKeeper        *paymentskeeper.Keeper
	FarmingKeeper         *farmingkeeper.Keeper
	PerpsKeeper           *perpskeeper.Keeper
	LendingKeeper         *lendingkeeper.Keeper
	PredictKeeper         *predictkeeper.Keeper
	LaunchpadKeeper       *launchpadkeeper.Keeper
	CLMMKeeper            *clmmkeeper.Keeper
	FlashLoanKeeper       *flashloankeeper.Keeper
	VaultKeeper           *vaultkeeper.Keeper
	PortfolioKeeper       *portfoliokeeper.Keeper
	OptionsKeeper         *optionskeeper.Keeper
	AIAgentKeeper         *aiagentkeeper.Keeper

	// Parallel execution engine
	ParallelExecutor *parallel.Executor

	mm           *module.Manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

func NewSyreenApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *SyreenApp {
	encodingConfig := MakeEncodingConfig()

	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(AppName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, minttypes.StoreKey,
		distrtypes.StoreKey, slashingtypes.StoreKey, govtypes.StoreKey, paramstypes.StoreKey,
		upgradetypes.StoreKey, consensusparamtypes.StoreKey, crisistypes.StoreKey,
		evidencetypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		ibcexported.StoreKey, ibctransfertypes.StoreKey,
		tokenfactorytypes.StoreKey, feemarkettypes.StoreKey, mevprotectiontypes.StoreKey,
		abstractaccounttypes.StoreKey, intenttypes.StoreKey, computetypes.StoreKey,
		evmtypes.StoreKey,
		dextypes.StoreKey,
		identitytypes.StoreKey,
		paymentstypes.StoreKey,
		farmingtypes.StoreKey,
		perpstypes.StoreKey,
		lendingtypes.StoreKey,
		flashloantypes.StoreKey,
		predicttypes.StoreKey,
		launchpadtypes.StoreKey,
		clmmtypes.StoreKey,
		vaulttypes.StoreKey,
		portfoliotypes.StoreKey,
		optionstypes.StoreKey,
		aiagenttypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys()

	app := &SyreenApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// ---------- Keepers ----------

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
	bApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(Bech32MainPrefix),
		Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		BlockedModuleAccountAddrs(maccPerms),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		stakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper.AddressCodec(),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights(appOpts),
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		DefaultNodeHome,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Feegrant
	app.FeegrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[feegrant.StoreKey]),
		app.AccountKeeper,
	)

	// Authz
	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	// Evidence
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		stakingKeeper,
		app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)

	// Staking hooks
	app.StakingKeeper = stakingKeeper
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// IBC Keeper (v10 — uses storeService, no stakingKeeper, no scopedKeeper)
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibcexported.StoreKey]),
		app.GetSubspace(ibcexported.ModuleName),
		app.UpgradeKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC Transfer Keeper (v10 — no portKeeper, no scopedKeeper, add msgRouter)
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibctransfertypes.StoreKey]),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC Router
	transferIBCModule := ibctransfer.NewIBCModule(app.TransferKeeper)
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// Token Factory
	app.TokenfactoryKeeper = tokenfactorykeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[tokenfactorytypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// DEX (AMM)
	dexK := dexkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[dextypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.TokenfactoryKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.DexKeeper = &dexK
	app.DexKeeper.SetStakingKeeper(app.StakingKeeper)

	// Fee Market
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[feemarkettypes.StoreKey]),
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// MEV Protection
	app.MEVProtectionKeeper = mevprotectionkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[mevprotectiontypes.StoreKey]),
		app.StakingKeeper,
		app.SlashingKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Abstract Account
	app.AbstractAccountKeeper = abstractaccountkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[abstractaccounttypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Intent System
	app.IntentKeeper = intentkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[intenttypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.IntentKeeper.SetDexKeeper(app.DexKeeper)
	app.IntentKeeper.SetTransferKeeper(&ibcTransferAdapter{keeper: app.TransferKeeper})

	// Liquidity Farming
	farmingK := farmingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[farmingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.FarmingKeeper = &farmingK

	// Perpetual Futures
	perpsK := perpskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[perpstypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	perpsK.SetDexKeeper(app.DexKeeper)
	app.PerpsKeeper = &perpsK

	// Lending & Borrowing
	lendingK := lendingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[lendingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	lendingK.SetDexKeeper(app.DexKeeper)
	app.LendingKeeper = &lendingK

	// Flash Loans
	flashloanK := flashloankeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[flashloantypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.FlashLoanKeeper = &flashloanK

	// Prediction Markets
	predictK := predictkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[predicttypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.PredictKeeper = &predictK

	// Token Launchpad (IDO)
	launchpadK := launchpadkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[launchpadtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	launchpadK.SetTokenFactoryKeeper(&app.TokenfactoryKeeper)
	launchpadK.SetDexKeeper(app.DexKeeper)
	app.LaunchpadKeeper = &launchpadK

	// Concentrated Liquidity (CLMM)
	clmmK := clmmkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[clmmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.CLMMKeeper = &clmmK

	// Vault Strategies (Yearn-style)
	vaultK := vaultkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[vaulttypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.VaultKeeper = &vaultK

	// Portfolio Dashboard + Social Feed + Competitions
	portfolioK := portfoliokeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[portfoliotypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.PortfolioKeeper = &portfolioK

	// On-Chain Options (Black-Scholes pricing, puts/calls)
	optionsK := optionskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[optionstypes.StoreKey]),
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	optionsK.SetDexKeeper(app.DexKeeper)
	app.OptionsKeeper = &optionsK

	// AI Agent Framework (autonomous on-chain trading agents)
	aiAgentK := aiagentkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[aiagenttypes.StoreKey]),
		app.BankKeeper,
		app.AccountKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	aiAgentK.SetDexKeeper(app.DexKeeper)
	app.AIAgentKeeper = &aiAgentK

	// Smart Contracts (Compute)
	app.ComputeKeeper = computekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[computetypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// EVM (Ethereum Virtual Machine)
	app.EVMKeeper = evmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Identity (KYC/Verification)
	app.IdentityKeeper = identitykeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[identitytypes.StoreKey]),
		app.AccountKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Payments (Invoicing & Settlement)
	app.PaymentsKeeper = paymentskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[paymentstypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Parallel Execution Engine
	app.ParallelExecutor = parallel.NewExecutor(logger, 0) // 0 = auto-detect CPU count

	// Governance (SDK v0.53: distrKeeper added as 6th argument)
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	govKeeper.SetLegacyRouter(govRouter)
	app.GovKeeper = *govKeeper.SetHooks(govtypes.NewMultiGovHooks())

	// ---------- Module Manager ----------

	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// Create the Tendermint light client module for IBC v10.
	// NewLightClientModule needs the codec and a StoreProvider (from IBC store service).
	storeProvider := ibcclienttypes.NewStoreProvider(runtime.NewKVStoreService(keys[ibcexported.StoreKey]))
	lightClientModule := ibctm.NewLightClientModule(appCodec, storeProvider)

	app.mm = module.NewManager(
		// Core SDK
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		// Additional SDK
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		// IBC (v10 — no capability module)
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibctm.NewAppModule(lightClientModule),
		// Custom
		tokenfactory.NewAppModule(app.TokenfactoryKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		mevprotection.NewAppModule(app.MEVProtectionKeeper),
		abstractaccount.NewAppModule(app.AbstractAccountKeeper),
		intent.NewAppModule(app.IntentKeeper),
		compute.NewAppModule(app.ComputeKeeper),
		evm.NewAppModule(app.EVMKeeper),
		dex.NewAppModule(app.DexKeeper),
		identity.NewAppModule(app.IdentityKeeper),
		payments.NewAppModule(app.PaymentsKeeper),
		farming.NewAppModule(app.FarmingKeeper),
		perps.NewAppModule(app.PerpsKeeper),
		lending.NewAppModule(app.LendingKeeper),
		flashloan.NewAppModule(app.FlashLoanKeeper),
		predict.NewAppModule(app.PredictKeeper),
		launchpad.NewAppModule(app.LaunchpadKeeper),
		clmm.NewAppModule(app.CLMMKeeper),
		vault.NewAppModule(app.VaultKeeper),
		portfolio.NewAppModule(app.PortfolioKeeper),
		options.NewAppModule(app.OptionsKeeper),
		aiagent.NewAppModule(app.AIAgentKeeper),
	)

	// Pre-block order (upgrade module must run before everything else)
	app.mm.SetOrderPreBlockers(upgradetypes.ModuleName)

	// Begin block order
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		ibctm.ModuleName,
		tokenfactorytypes.ModuleName,
		feemarkettypes.ModuleName,
		mevprotectiontypes.ModuleName,
		abstractaccounttypes.ModuleName,
		farmingtypes.ModuleName,  // distribute farming rewards
		perpstypes.ModuleName,    // funding rates + liquidations
		lendingtypes.ModuleName,  // accrue lending interest
		predicttypes.ModuleName,  // auto-close expired prediction markets
		launchpadtypes.ModuleName, // auto-finalize expired launches
		clmmtypes.ModuleName,
		flashloantypes.ModuleName,
		vaulttypes.ModuleName,    // auto-compound vault yields
		portfoliotypes.ModuleName,
		optionstypes.ModuleName,  // expire options
		aiagenttypes.ModuleName,  // execute AI agent strategies
		dextypes.ModuleName,      // risk checks run before intent trading execution
		intenttypes.ModuleName,
		computetypes.ModuleName,
		evmtypes.ModuleName,
		identitytypes.ModuleName,
		paymentstypes.ModuleName,
	)

	// End block order (capability removed)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		ibctm.ModuleName,
		tokenfactorytypes.ModuleName,
		feemarkettypes.ModuleName,
		mevprotectiontypes.ModuleName,
		abstractaccounttypes.ModuleName,
		intenttypes.ModuleName,
		computetypes.ModuleName,
		evmtypes.ModuleName,
		dextypes.ModuleName,
		farmingtypes.ModuleName,
		perpstypes.ModuleName,
		lendingtypes.ModuleName,
		predicttypes.ModuleName,
		launchpadtypes.ModuleName,
		clmmtypes.ModuleName,
		flashloantypes.ModuleName,
		vaulttypes.ModuleName,
		portfoliotypes.ModuleName,
		optionstypes.ModuleName,
		aiagenttypes.ModuleName,
		identitytypes.ModuleName,
		paymentstypes.ModuleName,
	)

	// Genesis order (capability removed)
	genesisModuleOrder := []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		ibctm.ModuleName,
		tokenfactorytypes.ModuleName,
		feemarkettypes.ModuleName,
		mevprotectiontypes.ModuleName,
		abstractaccounttypes.ModuleName,
		intenttypes.ModuleName,
		computetypes.ModuleName,
		evmtypes.ModuleName,
		dextypes.ModuleName,
		farmingtypes.ModuleName,
		perpstypes.ModuleName,
		lendingtypes.ModuleName,
		predicttypes.ModuleName,
		launchpadtypes.ModuleName,
		clmmtypes.ModuleName,
		flashloantypes.ModuleName,
		vaulttypes.ModuleName,
		portfoliotypes.ModuleName,
		optionstypes.ModuleName,
		aiagenttypes.ModuleName,
		identitytypes.ModuleName,
		paymentstypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	if err := app.mm.RegisterServices(app.configurator); err != nil {
		panic(err)
	}

	// Initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// Set handlers
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: authante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				FeegrantKeeper:  app.FeegrantKeeper,
				SigGasConsumer:  authante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:                 app.IBCKeeper,
			FeeMarketKeeper:           app.FeeMarketKeeper,
			FeeMarketBankKeeper:       app.BankKeeper,
			AbstractAccountKeeper:     app.AbstractAccountKeeper,
			AbstractAccountBankKeeper: app.BankKeeper,
			EVMKeeper:                 app.EVMKeeper,
			FullAccountKeeper:         app.AccountKeeper,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}
	app.SetAnteHandler(anteHandler)

	// H-13: PostHandler to record session key spend usage after successful tx execution.
	// This ensures spend limits are only consumed when the transaction actually succeeds.
	postHandler := sdk.ChainPostDecorators(
		aaante.NewSessionKeyPostDecorator(app.AbstractAccountKeeper),
	)
	app.SetPostHandler(postHandler)

	// MEV protection: register ABCI proposal handlers for fair ordering enforcement
	app.SetPrepareProposal(app.PrepareProposalHandler())
	app.SetProcessProposal(app.ProcessProposalHandler())

	// Register upgrade handlers only for the current binary version.
	// The genesis/pre-upgrade binary must NOT register handlers for future upgrades,
	// or the PreBlocker will panic with "BINARY UPDATED BEFORE TRIGGER".
	appVersion := version.Version
	if appVersion == "1.1.0" {
		app.UpgradeKeeper.SetUpgradeHandler("v1.1.0", func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("running v1.1.0 upgrade handler", "height", sdkCtx.BlockHeight())
			// Supply fix (410M→400M) applied in genesis. No burn needed.
			// Voting period (10min) also set in genesis. No param change needed.
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		})
	}

	// v2.0.0: Founder vesting schedule (6-month cliff + 24-month linear vest)
	app.UpgradeKeeper.SetUpgradeHandler("v2.0.0",
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			logger.Info("applying v2.0.0 upgrade: founder vesting schedule")

			founderAddr, err := sdk.AccAddressFromBech32("syreen1rrj8dx99djxkjx4rlfy2xryca9g7fujjjfclvz")
			if err != nil {
				return nil, err
			}

			acct := app.AccountKeeper.GetAccount(sdkCtx, founderAddr)
			if acct == nil {
				logger.Info("founder account not found, skipping")
				return app.mm.RunMigrations(ctx, app.configurator, fromVM)
			}

			baseAcct, ok := acct.(*authtypes.BaseAccount)
			if !ok {
				logger.Info("founder account not a BaseAccount, skipping")
				return app.mm.RunMigrations(ctx, app.configurator, fromVM)
			}

			now := sdkCtx.BlockTime()
			cliffEnd := now.Add(6 * 30 * 24 * time.Hour) // ~6 months
			founderBalance := app.BankKeeper.GetBalance(sdkCtx, founderAddr, "usyreen")
			totalVesting := founderBalance.Amount

			monthlyAmount := totalVesting.Quo(math.NewInt(24))
			remainder := totalVesting.Sub(monthlyAmount.Mul(math.NewInt(24)))

			cliffSeconds := int64(cliffEnd.Sub(now).Seconds())
			monthSeconds := int64(30 * 24 * 3600)

			periods := make([]vestingtypes.Period, 0, 24)
			for i := 0; i < 24; i++ {
				amt := monthlyAmount
				if i == 23 {
					amt = amt.Add(remainder)
				}
				length := monthSeconds
				if i == 0 {
					length = cliffSeconds + monthSeconds
				}
				periods = append(periods, vestingtypes.Period{
					Length: length,
					Amount: sdk.NewCoins(sdk.NewCoin("usyreen", amt)),
				})
			}

			vestingAcct, err := vestingtypes.NewPeriodicVestingAccount(
				baseAcct,
				sdk.NewCoins(sdk.NewCoin("usyreen", totalVesting)),
				now.Unix(),
				periods,
			)
			if err != nil {
				return nil, err
			}
			app.AccountKeeper.SetAccount(sdkCtx, vestingAcct)

			logger.Info("founder vesting account created",
				"address", founderAddr.String(),
				"total_vesting", totalVesting.String(),
				"cliff_end", cliffEnd.String(),
				"fully_vested", cliffEnd.Add(24*30*24*time.Hour).String(),
			)

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			// If latest version is corrupt (e.g. mid-commit crash), try
			// rolling back one version which should be fully committed.
			app.Logger().Error("LoadLatestVersion failed, attempting rollback", "error", err)
			cms := app.CommitMultiStore()
			latest := cms.LatestVersion()
			if latest > 1 {
				app.Logger().Info("Attempting to load previous version", "version", latest-1)
				if err2 := cms.LoadVersion(latest - 1); err2 != nil {
					panic(fmt.Errorf("error loading latest version: %w; rollback to %d also failed: %w", err, latest-1, err2))
				}
				app.Logger().Info("Successfully rolled back to previous version", "version", latest-1)
			} else {
				panic(fmt.Errorf("error loading latest version: %w", err))
			}
		}

		// Wrap the commit multistore in FallbackQueryMultiStore so that
		// CacheMultiStoreWithVersion falls back to latest state instead of
		// failing with "version does not exist" (SDK v0.50.x IAVL bug).
		// This fixes ALL standard gRPC/REST queries and tx simulation.
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/13317
		app.SetQueryMultiStore(NewFallbackQueryMultiStore(app.CommitMultiStore()))
	}

	return app
}

func (app *SyreenApp) Name() string { return app.BaseApp.Name() }

func (app *SyreenApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

func (app *SyreenApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

func (app *SyreenApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

func (app *SyreenApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}
	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		return nil, err
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *SyreenApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *SyreenApp) LegacyAmino() *codec.LegacyAmino   { return app.legacyAmino }
func (app *SyreenApp) AppCodec() codec.Codec               { return app.appCodec }
func (app *SyreenApp) InterfaceRegistry() types.InterfaceRegistry { return app.interfaceRegistry }
func (app *SyreenApp) TxConfig() client.TxConfig            { return app.txConfig }

func (app *SyreenApp) GetKey(storeKey string) *storetypes.KVStoreKey       { return app.keys[storeKey] }
func (app *SyreenApp) GetTKey(storeKey string) *storetypes.TransientStoreKey { return app.tkeys[storeKey] }
func (app *SyreenApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey  { return app.memKeys[storeKey] }

func (app *SyreenApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// ibcTransferAdapter wraps the IBC transfer keeper to satisfy the intent module's TransferKeeper interface.
type ibcTransferAdapter struct {
	keeper ibctransferkeeper.Keeper
}

func (a *ibcTransferAdapter) Transfer(ctx context.Context, msg *intenttypes.IBCTransferMsg) (*intenttypes.IBCTransferResponse, error) {
	ibcMsg := ibctransfertypes.NewMsgTransfer(
		msg.SourcePort,
		msg.SourceChannel,
		msg.Token,
		msg.Sender,
		msg.Receiver,
		ibcclienttypes.NewHeight(0, msg.TimeoutHeight),
		msg.TimeoutTimestamp,
		"", // memo
	)

	resp, err := a.keeper.Transfer(ctx, ibcMsg)
	if err != nil {
		return nil, err
	}

	return &intenttypes.IBCTransferResponse{Sequence: resp.Sequence}, nil
}

func (app *SyreenApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register custom direct-state query endpoints that bypass the broken
	// CacheMultiStoreWithVersion in SDK v0.50.x IAVL store integration.
	app.RegisterCustomQueryRoutes(apiSvr.Router, clientCtx)

	// Register standard Cosmos SDK REST endpoints (proxied through keepers)
	// so that Keplr, Mintscan, and other ecosystem tools work out of the box.
	app.RegisterStandardCosmosRoutes(apiSvr.Router)

	// Start Ethereum JSON-RPC server for MetaMask/EVM wallet compatibility.
	// C5: Bind to loopback only — the public-facing EVM RPC must be served via
	// the nginx reverse proxy on the same host (which terminates TLS and applies
	// rate limiting / auth). Binding 0.0.0.0 here would expose the JSON-RPC
	// directly to the internet, bypassing the nginx layer that fronts every
	// other endpoint (LCD on 1317, CometBFT RPC on 26657, gRPC, etc.).
	_, err := jsonrpc.StartJSONRPCServer(app, "127.0.0.1:8545", "http://127.0.0.1:26657", app.Logger())
	if err != nil {
		app.Logger().Error("Failed to start JSON-RPC server", "error", err)
	}
}

func (app *SyreenApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *SyreenApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

func (app *SyreenApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

func (app *SyreenApp) SimulationManager() *module.SimulationManager { return app.sm }

func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

func BlockedModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	// The intent module MUST be able to receive swap outputs from the DEX module
	// during DCA/limit order execution. If it stays blocked, SendCoinsFromModuleToAccount
	// from dex -> intent fails inside dexKeeper.Swap(), leaving the input stuck in the
	// DEX module bank account with no output delivered.
	delete(modAccAddrs, authtypes.NewModuleAddress(intenttypes.ModuleName).String())
	return modAccAddrs
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	return paramsKeeper
}

func skipUpgradeHeights(appOpts servertypes.AppOptions) map[int64]bool {
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	return skipUpgradeHeights
}

func (app *SyreenApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	ctx := app.NewContext(true)
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
	}

	genState, err := app.mm.ExportGenesis(ctx, app.appCodec)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// Commit overrides BaseApp.Commit to update the query store snapshot
// after each block commit. This ensures gRPC/REST queries always read
// from the last COMMITTED state, not the in-progress FinalizeBlock state.
func (app *SyreenApp) Commit() (*abci.ResponseCommit, error) {
	resp, err := app.BaseApp.Commit()
	if err != nil {
		return resp, err
	}
	// After commit, update the fallback query store to point to the
	// freshly committed state.
	app.SetQueryMultiStore(NewFallbackQueryMultiStore(app.CommitMultiStore()))
	return resp, err
}

// --- JSON-RPC AppStateProvider interface ---

// GetCommitMultiStore implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetCommitMultiStore() storetypes.CommitMultiStore {
	return app.CommitMultiStore()
}

// GetBankKeeper implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

// GetEVMKeeper implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetEVMKeeper() evmkeeper.Keeper {
	return app.EVMKeeper
}

// GetLogger implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetLogger() log.Logger {
	return app.Logger()
}

// GetTxConfig implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetTxConfig() client.TxConfig {
	return app.txConfig
}

// GetAccountKeeper implements jsonrpc.AppStateProvider
func (app *SyreenApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}
