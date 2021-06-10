package app

import (
	"encoding/json"
	"github.com/anathatech/project-anatha/x/evidence"
	appConfig "github.com/anathatech/project-anatha/config"
	"github.com/anathatech/project-anatha/x/distribution"
	"github.com/anathatech/project-anatha/x/fee"
	feeclient "github.com/anathatech/project-anatha/x/fee/client"
	"github.com/anathatech/project-anatha/x/hra"
	hraclient "github.com/anathatech/project-anatha/x/hra/client"
	"github.com/anathatech/project-anatha/x/treasury"
	"io"
	"os"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/anathatech/project-anatha/x/genutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	distributionclient "github.com/anathatech/project-anatha/x/distribution/client"
	treasuryclient "github.com/anathatech/project-anatha/x/treasury/client"
	gov "github.com/anathatech/project-anatha/x/governance"
	"github.com/anathatech/project-anatha/x/mint"
	"github.com/anathatech/project-anatha/x/slashing"
	"github.com/anathatech/project-anatha/x/staking"
	"github.com/anathatech/project-anatha/x/upgrade"
	upgradeclient "github.com/anathatech/project-anatha/x/upgrade/client"
)


var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.anathacli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.anathad")

	// NewBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		gov.NewAppModuleBasic(
			upgradeclient.ProposalHandler,
			hraclient.RegisterBlockchainIdProposalHandler,
			hraclient.RemoveBlockchainIdProposalHandler,
			distributionclient.DevelopmentFundDistributionProposalHandler,
			distributionclient.SecurityTokenFundDistributionProposalHandler,
			feeclient.AddFeeExcludedMessageProposalHandler,
			feeclient.RemoveFeeExcludedMessageProposalHandler,
			treasuryclient.AddBuyBackLiquidityProposalHandler,
			treasuryclient.RemoveBuyBackLiquidityProposalHandler,
			treasuryclient.BurnDistributionProfitsProposalHandler,
			treasuryclient.TransferFromDistributionProfitsToBuyBackLiquidityProposalHandler,
			treasuryclient.TransferFromTreasuryToSwapEscrowProposalHandler,
			treasuryclient.TransferFromSwapEscrowToBuyBackProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		hra.AppModuleBasic{},
		treasury.AppModuleBasic{},
		distribution.AppModuleBasic{},
		fee.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:           nil,
		mint.ModuleName:                 {supply.Minter},
		treasury.ModuleName:             {supply.Minter},
		treasury.BuyBackLiquidityFundModuleName: {supply.Minter, supply.Burner},
		treasury.BuyBackFundModuleName:           nil,
		treasury.DistributionProfitsModuleName:   {supply.Burner},
		treasury.TreasuryEscrowModuleName:        nil,
		treasury.SwapEscrowModuleName:            nil,
		staking.BondedPoolName:                   {supply.Burner, supply.Staking},
		staking.NotBondedPoolName:                {supply.Burner, supply.Staking},
		gov.ModuleName:                           nil,
		distribution.AmcModuleName:               nil,
		distribution.NvrpModuleName:              nil,
		distribution.NvrpDistributionModuleName:  nil,
		distribution.HRAHolderRewardModuleName:   nil,
		distribution.DevelopmentFundModuleName:   nil,
		distribution.SecurityTokenFundModuleName: nil,
		distribution.SavingsModuleName:           nil,
		distribution.SavingsDistributionModuleName: nil,
	}
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc.Seal()
}

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(appConfig.Bech32MainPrefix, appConfig.Bech32MainPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(appConfig.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, appConfig.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(appConfig.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, appConfig.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
}

type AnathaApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// Keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	upgradeKeeper  upgrade.Keeper
	evidenceKeeper evidence.Keeper
	treasuryKeeper treasury.Keeper
	distributionKeeper distribution.Keeper
	feeKeeper 		fee.Keeper

	hraKeeper 		hra.Keeper

	// Module Manager
	mm *module.Manager
}

// NewAnathaApp is a constructor function for AnathaApp
func NewAnathaApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, baseAppOptions ...func(*bam.BaseApp)) *AnathaApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appConfig.AppName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		mint.StoreKey,
		slashing.StoreKey,
		gov.StoreKey,
		params.StoreKey,
		evidence.StoreKey,
		upgrade.StoreKey,
		hra.StoreKey,
		treasury.StoreKey,
		distribution.StoreKey,
		fee.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &AnathaApp{
		BaseApp: 		bApp,
		cdc:     		cdc,
		invCheckPeriod: invCheckPeriod,
		keys:    		keys,
		tkeys:   		tkeys,
		subspaces: 		make(map[string]params.Subspace),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	// Set specific supspaces
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[hra.ModuleName] = app.paramsKeeper.Subspace(hra.DefaultParamspace)
	app.subspaces[treasury.ModuleName] = app.paramsKeeper.Subspace(treasury.DefaultParamspace)
	app.subspaces[distribution.ModuleName] = app.paramsKeeper.Subspace(distribution.DefaultParamspace)
	app.subspaces[fee.ModuleName] = app.paramsKeeper.Subspace(fee.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	hraKeeper := hra.NewKeeper(
		app.bankKeeper,
		app.accountKeeper,
		app.supplyKeeper,
		app.cdc,
		keys[hra.StoreKey],
		app.subspaces[hra.ModuleName],
		distribution.AmcModuleName,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
		&hraKeeper,
	)

	app.mintKeeper = mint.NewKeeper(
		app.cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName],
		app.supplyKeeper, distribution.AmcModuleName, distribution.NvrpModuleName,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		app.subspaces[slashing.ModuleName],
	)

	app.crisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.supplyKeeper, distribution.AmcModuleName,
	)

	app.feeKeeper = fee.NewKeeper(
		app.cdc,
		keys[fee.StoreKey],
		app.subspaces[fee.ModuleName],
	)

	app.treasuryKeeper = treasury.NewKeeper(
		app.cdc,
		keys[treasury.StoreKey],
		app.subspaces[treasury.ModuleName],
		app.supplyKeeper,
		app.accountKeeper,
		app.bankKeeper,
	)

	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.cdc)

	app.upgradeKeeper.SetUpgradeHandler("swap", func(ctx sdk.Context, plan upgrade.Plan) {
		// Add "treasury/swap" to fee excluded messages
		app.feeKeeper.SetFeeExcludedMessage(ctx, "treasury/swap")

		// Transfer swap module balance from initial distribution account
		// TODO: Change the address if necessary
		swapAccountAddress := "anatha1atzrzll4k9r7et8pa6dx70vc3x8s7wkww0v7gf"

		address, _ := sdk.AccAddressFromBech32(swapAccountAddress)

		// TODO: Change the amount if necessary
		amount := sdk.NewCoins(sdk.NewInt64Coin(appConfig.DefaultDenom, 30000000000000000))

		_ = app.supplyKeeper.SendCoinsFromAccountToModule(ctx, address, treasury.SwapEscrowModuleName, amount)

		// Remove devnet minted buyback liquidity funds
		app.treasuryKeeper.BurnCoinsFromBuyBackLiquidityFund(ctx, sdk.NewCoins(sdk.NewInt64Coin(appConfig.DefaultStableDenom, 10000000000000)))

		app.feeKeeper.SetParams(
			ctx,
			fee.NewParams(
				app.feeKeeper.FeePercentage(ctx),
				app.feeKeeper.MinimumFee(ctx),
				sdk.NewCoins(sdk.NewInt64Coin(appConfig.DefaultDenom, 100000000)),
			),
		)

		// Update risk assesment duration to 24 hours
		treasuryParams := app.treasuryKeeper.GetParams(ctx)
		treasuryParams.RiskAssessmentDuration = time.Hour * 24
		app.treasuryKeeper.SetParams(ctx, treasuryParams)

	})

	app.upgradeKeeper.SetUpgradeHandler("delay", func(ctx sdk.Context, plan upgrade.Plan) {
		// Update risk assesment duration to 1 hour
		treasuryParams := app.treasuryKeeper.GetParams(ctx)
		treasuryParams.RiskAssessmentDuration = time.Hour
		app.treasuryKeeper.SetParams(ctx, treasuryParams)
	})

	// create evidence keeper with evidence router
	evidenceKeeper := evidence.NewKeeper(
		app.cdc, keys[evidence.StoreKey], app.subspaces[evidence.ModuleName], &stakingKeeper, app.slashingKeeper,
	)

	evidenceRouter := evidence.NewRouter()
	evidenceKeeper.SetRouter(evidenceRouter)

	app.evidenceKeeper = *evidenceKeeper

	app.distributionKeeper = distribution.NewKeeper(
		app.cdc,
		keys[distribution.StoreKey],
		app.subspaces[distribution.ModuleName],
		app.supplyKeeper,
		&stakingKeeper,
		&hraKeeper,
	)

	app.hraKeeper = *hraKeeper.SetHooks(
		hra.NewMultiNameHooks(app.distributionKeeper.NameHooks()),
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		app.subspaces[gov.ModuleName],
		govRouter,
		app.accountKeeper,
	)

	govRouter.AddRoute(gov.RouterKey, gov.NewGovernanceProposalHandler(app.govKeeper)).
		// AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)). // 328 - THIS - Controls Governance of Parameters
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(hra.RouterKey, hra.NewGovernanceProposalHandler(app.hraKeeper)).
		AddRoute(fee.RouterKey, fee.NewGovernanceProposalHandler(app.feeKeeper)).
		AddRoute(distribution.RouterKey, distribution.NewDistributionProposalHandler(app.distributionKeeper)).
		AddRoute(treasury.RouterKey, treasury.NewTreasuryProposalHandler(app.treasuryKeeper))

	govRouter.Seal()

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distributionKeeper.StakingHooks(),
			app.slashingKeeper.Hooks()),
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		hra.NewAppModule(app.hraKeeper, app.bankKeeper, app.accountKeeper, app.supplyKeeper),
		treasury.NewAppModule(app.treasuryKeeper, app.accountKeeper, app.supplyKeeper, app.bankKeeper),
		distribution.NewAppModule(app.distributionKeeper, app.supplyKeeper),
		fee.NewAppModule(app.feeKeeper),
	)

	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, mint.ModuleName, distribution.ModuleName, treasury.ModuleName, slashing.ModuleName, staking.ModuleName)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName, hra.ModuleName, treasury.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		hra.ModuleName,
		supply.ModuleName,
		crisis.ModuleName,
		evidence.ModuleName,
		treasury.ModuleName,
		fee.ModuleName,
		genutil.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		app.NewAnteHandler(),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *AnathaApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *AnathaApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *AnathaApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *AnathaApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *AnathaApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *AnathaApp) Codec() *codec.Codec {
	return app.cdc
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

func init() {
	backupGenesisDefaults()
	setGenesisDefaults()
	RegisterCoinDenoms()
}

var stakingDefaultGenesisState staking.GenesisState
var crisisDefaultGenesisState crisis.GenesisState

func backupGenesisDefaults() {
	stakingDefaultGenesisState = staking.DefaultGenesisState()
	crisisDefaultGenesisState = crisis.DefaultGenesisState()
}

func setGenesisDefaults() {
	// Override module defaults for use in testnets and the default init functionality.
	staking.DefaultGenesisState = func () staking.GenesisState {
		stakingDefaultGenesisState.Params.BondDenom = appConfig.DefaultDenom
		return stakingDefaultGenesisState
	}
	crisis.DefaultGenesisState = func () crisis.GenesisState {
		crisisDefaultGenesisState.ConstantFee.Denom = appConfig.DefaultDenom
		crisisDefaultGenesisState.ConstantFee.Amount = sdk.NewInt(50000000000000)
		return crisisDefaultGenesisState
	}
}

func RegisterCoinDenoms() {
	_ = sdk.RegisterDenom("anatha", sdk.OneDec())
	_ = sdk.RegisterDenom("sense", sdk.NewDecWithPrec(1, 3))
	_ = sdk.RegisterDenom("pin", sdk.NewDecWithPrec(1, 8))

	_ = sdk.RegisterDenom("usd", sdk.OneDec())
	_ = sdk.RegisterDenom("din", sdk.NewDecWithPrec(1, 10))
}